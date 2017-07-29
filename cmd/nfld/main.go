package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"log/syslog"
	"os"

	"github.com/ameske/nfl-pickem/http"
	"github.com/ameske/nfl-pickem/sqlite3"
)

// For now, we will let all of these things be global since it's easier
var (
	slog *syslog.Writer
)

type config struct {
	Server struct {
		AuthKey    string `json:"authKey"`
		EncryptKey string `json:"encryptKey"`
		Database   string `json:"databaseFile"`
		LogosDir   string `json:"logosDirectory"`
	} `json:"server"`
	Email struct {
		Enabled     bool   `json:"enabled"`
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

func main() {
	configFile := flag.String("config", "/opt/ameske/gonfl/conf.json", "Path to server config file")
	debug := flag.Bool("debug", false, "used when running the server out of the source repo")
	dbFile := flag.String("db", "", "override the configuration database")
	/*
		runUpdate := flag.String("update", "", "update scores and results using given results JSON file")
		gradeOnly := flag.Bool("grade", false, "don't run in daemon mode and just grade the given year and week")
		year := flag.Int("year", -1, "year for batch processing mode")
		week := flag.Int("week", -1, "week for batch processing mode")
	*/

	flag.Parse()

	var err error
	slog, err = syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL0, "nfl-pickem")
	if err != nil {
		log.Fatal("Could not connect to syslog:", err)
	}

	var c config

	if !*debug {
		c = loadConfig(*configFile)
		log.SetOutput(slog)
	} else {
		log.SetOutput(io.MultiWriter(slog, os.Stdout))
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		c.Server.AuthKey = base64.StdEncoding.EncodeToString([]byte("something secret"))
		c.Server.EncryptKey = base64.StdEncoding.EncodeToString([]byte("something secret"))
		c.Server.Database = "nfl.db"
		c.Server.LogosDir = "/Users/ameske/Documents/go/src/github.com/ameske/nfl-pickem/logos"
		c.Email.Enabled = false
	}

	if *dbFile != "" {
		c.Server.Database = *dbFile
	}

	db, err := sqlite3.NewDatastore(c.Server.Database)
	if err != nil {
		log.Fatal(err)
	}

	/*
		var n Notifier
		if c.Email.Enabled {
			n, err = NewEmailNotifier(c.Email.SMTPAddress, c.Email.Sender, c.Email.Password)
			if err != nil {
				n = nullNotifier{}
			}

		} else if *debug {
			n = fsNotifier{}
		} else {

			n = nullNotifier{}
		}

	*/

	scheduleUpdates(db)

	hashKey, encryptKey, err := parseSecureCookieKeys(c.Server.AuthKey, c.Server.EncryptKey)
	if err != nil {
		log.Fatal(err)
	}

	server, err := http.NewServer("0.0.0.0:61389", hashKey, encryptKey, db)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Start())
}
