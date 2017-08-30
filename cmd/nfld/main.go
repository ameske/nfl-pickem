package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"os"
	"time"

	nflpickem "github.com/ameske/nfl-pickem"
	"github.com/ameske/nfl-pickem/http"
	"github.com/ameske/nfl-pickem/sqlite3"
)

// For now, we will let all of these things be global since it's easier

type config struct {
	Server struct {
		AuthKey    string `json:"authKey"`
		EncryptKey string `json:"encryptKey"`
		Database   string `json:"databaseFile"`
		LogosDir   string `json:"logosDirectory"`
		Autoupdate bool   `json:"autoupdateEnabled"`
	} `json:"server"`
	Email struct {
		Enabled     bool   `json:"enabled"`
		Type        string `json:"type"`
		Sender      string `json:"sendAsAddress"`
		Password    string `json:"password"`
		SMTPAddress string `json:"smtpAddress"`
	} `json:"email"`
}

func loadConfig(path string) config {
	configBytes, err := ioutil.ReadFile(path)

	config := config{}
	err = json.Unmarshal(configBytes, &config)

	if err != nil {
		log.Fatal(err)
	}

	return config
}

func parseSecureCookieKeys(b64HashKey, b64EncryptKey string) (hashKey []byte, encryptKey []byte, err error) {
	hashKey, err = base64.StdEncoding.DecodeString(b64HashKey)
	if err != nil {
		return nil, nil, err
	}

	encryptKey, err = base64.StdEncoding.DecodeString(b64EncryptKey)
	if err != nil {
		return nil, nil, err
	}

	return
}

func setupNotifier(c config) (n nflpickem.Notifier, err error) {
	if !c.Email.Enabled {
		return nullNotifier{}, nil
	}

	switch c.Email.Type {
	case "fs":
		n, err = fsNotifier{}, nil
	case "email":
		n, err = NewEmailNotifier(c.Email.SMTPAddress, c.Email.Sender, c.Email.Password)
	default:
		n, err = nil, fmt.Errorf("unrecognized e-mail type: %s", c.Email.Type)
	}

	return n, err
}

func main() {
	var dbFile, configFile, timeStr string
	var stdout bool

	flag.StringVar(&configFile, "config", "/opt/ameske/gonfl/conf.json", "Path to server config file")
	flag.StringVar(&dbFile, "db", "", "override the configuration file's databsae location")
	flag.BoolVar(&stdout, "stdout", false, "log to the console instead of syslog")
	flag.StringVar(&timeStr, "time", "", "override the current internal time of the server (use 'unix date' format)")

	flag.Parse()

	var err error

	var c config
	c = loadConfig(configFile)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if stdout {
		log.SetOutput(os.Stdout)
	} else {
		slog, err := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL0, "nfl-pickem")
		if err != nil {
			log.Fatal("Could not connect to syslog:", err)
		}
		log.SetOutput(slog)
	}

	if dbFile != "" {
		c.Server.Database = dbFile
	}

	db, err := sqlite3.NewDatastore(c.Server.Database)
	if err != nil {
		log.Fatal(err)
	}

	notifier, err := setupNotifier(c)
	if err != nil {
		log.Fatal(err)
	}

	if c.Server.Autoupdate {
		scheduleUpdates(db)
	}

	hashKey, encryptKey, err := parseSecureCookieKeys(c.Server.AuthKey, c.Server.EncryptKey)
	if err != nil {
		log.Fatal(err)
	}

	var timeSource http.TimeSource
	if timeStr != "" {
		t, err := time.Parse(time.UnixDate, timeStr)
		if err != nil {
			log.Fatal(err)
		}

		timeSource = NewCustomTime(t)
	} else {
		timeSource = http.DefaultTimesource
	}

	prefix := "/api"
	server, err := http.NewServer("0.0.0.0:61389", prefix, hashKey, encryptKey, db, notifier, timeSource)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Start())
}

// customTime implements the TimeSource interface provided by package HTTP
type customTime struct {
	now time.Time
}

// NewCustomTime constructs a time source that always uses the provided time as "now"
func NewCustomTime(t time.Time) customTime {
	return customTime{
		now: t,
	}
}

// Now returns the static time the customTime was created with
func (ct customTime) Now() time.Time {
	return ct.now
}
