package main

import (
	"log"
	"net/http"

	"github.com/ameske/nfl-pickem/jsonhttp"
	"github.com/gorilla/sessions"
)

type CredentialChecker interface {
	CheckCredentials(username string, password string) (bool, error)
	IsAdmin(usernamd string) (bool, error)
}

// Login processes a login request and sets a cookie with a session handle on success
func Login(db CredentialChecker, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		// Attempt login, taking the user back to the login page with an error message if failed
		u := r.FormValue("username")
		p := r.FormValue("password")
		ok, err := db.CheckCredentials(u, p)
		if err != nil {
			log.Println(err)
			jsonhttp.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		} else if !ok {
			jsonhttp.WriteError(w, http.StatusUnauthorized, "invalid username/password")
			return
		}

		// Set session information
		session, _ := store.Get(r, "LoginState")
		session.Values["status"] = "loggedin"
		session.Values["user"] = u
		session.Values["admin"], err = db.IsAdmin(u)
		if err != nil {
			log.Println(err)
			jsonhttp.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		err = session.Save(r, w)
		if err != nil {
			log.Println(err)
			jsonhttp.WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		jsonhttp.WriteSuccess(w, "successfully logged in")
	}
}

// Logout clears the session information, which effectively logs the user out.
func Logout(store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "LoginState")
		session.Options.MaxAge = -1
		session.Save(r, w)

		jsonhttp.WriteSuccess(w, "successful logout")
	}
}
