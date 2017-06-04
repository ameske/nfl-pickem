package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/ameske/nfl-pickem/api"
	"github.com/ameske/nfl-pickem/jsonhttp"
	"github.com/ameske/nfl-pickem/sqlite3"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/husobee/vestigo"
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

func AddUserInfo(h http.HandlerFunc, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "LoginState")
		if session.Values["status"] == "loggedin" {
			context.Set(r, "user", session.Values["user"].(string))
			context.Set(r, "admin", session.Values["admin"].(bool))
		}

		h(w, r)
	}
}

func RequireLogin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := context.Get(r, "user").(string)
		if ok && u != "" {
			jsonhttp.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
	}
}

func AdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := context.GetOk(r, "user")
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		a, ok := context.GetOk(r, "admin")
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		user, ok := u.(string)
		if !ok {
			log.Println("Could not cast user to string when checking admin only")
			http.Error(w, "", http.StatusInternalServerError)
		}

		isAdmin, ok := a.(bool)
		if !ok {
			log.Println("Could not cast admin to bool when checking admin only")
			http.Error(w, "", http.StatusInternalServerError)
		}

		if user == "" || !isAdmin {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		h(w, r)
	}
}

func configureSessionStore(b64HashKey, b64EncryptKey string) (sessions.Store, error) {
	hk, err := base64.StdEncoding.DecodeString(b64HashKey)
	if err != nil {
		return nil, err
	}

	ek, err := base64.StdEncoding.DecodeString(b64EncryptKey)
	if err != nil {
		return nil, err
	}

	return sessions.NewCookieStore(hk, ek), nil
}

func currentUser(r *http.Request) (string, bool) {
	user := ""
	u, ok := context.GetOk(r, "user")
	if ok {
		user, ok = u.(string)
	}

	admin := false
	a, ok := context.GetOk(r, "admin")
	if ok {
		admin, ok = a.(bool)
	}

	return user, admin
}

func yearWeek(r *http.Request) (int, int) {
	v := mux.Vars(r)
	y, _ := strconv.ParseInt(v["year"], 10, 32)
	w, _ := strconv.ParseInt(v["week"], 10, 32)
	return int(y), int(w)
}

func dumpRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r, false)
		if err != nil {
			fmt.Println("Could not dump request:", err)
		}

		fmt.Println(string(b))

		h(w, r)
	}
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
	// var store sessions.Store

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

		store, err = configureSessionStore(c.Server.AuthKey, c.Server.EncryptKey)
		if err != nil {
			log.Fatal(err)
		}

		scheduleUpdates()
	*/

	router := vestigo.NewRouter()

	router.Get("/current", api.CurrentWeek(db))
	router.Get("/games", api.Games(db))
	router.Get("/picks", api.GetPicks(db))
	router.Get("/results", api.Results(db))
	router.Get("/totals", api.WeeklyTotals(db))

	log.Printf("NFL Pick-Em Pool listening on port 61389")
	log.Fatal(http.ListenAndServe("0.0.0.0:61389", router))
}
