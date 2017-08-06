package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ameske/nfl-pickem"
)

// pickManager is the interface that defines the ability to retrieve and make picks
type pickManager interface {
	nflpickem.PickRetriever
	nflpickem.Picker
}

// picks retrieves a user's picks for the provided week of the NFL season, OR updates
// the picks based on provided picks in the request body.
//
// URL Parameters:
//	year: Specifies the current year, Required
//	week: Specifies the current week, Required
func picks(db pickManager, notifier nflpickem.Notifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := retrieveUser(r.Context())
		if err == errNoUser {
			WriteJSONError(w, http.StatusUnauthorized, "login required")
			return
		} else if err != nil {
			WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if r.Method == "GET" {
			getPicks(user, db, w, r)
		} else if r.Method == "POST" {
			postPicks(user, db, notifier, w, r)
		} else {
			WriteJSONError(w, http.StatusMethodNotAllowed, "only GET or POST allowed")
		}
	}
}

// GetPicks returns the set of picks for the given user, year, and week.
func getPicks(user nflpickem.User, db nflpickem.PickRetriever, w http.ResponseWriter, r *http.Request) {
	yearStr := r.FormValue("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "year query parameter must be integer")
		return
	}

	weekStr := r.FormValue("week")
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "week query parameter must be integer")
		return
	}

	username := r.FormValue("username")

	picks, err := db.UserPicks(username, year, week)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, picks)
}

// MakePicks processes an array of JSON representation of pick selections.
//
// In the event duplicate picks for the same game are made,
// the last pick is always the pick that is stored.
//
// This endpoint restricts the set of picks to be for a pre-declared user,
// declared in the URL.
//
// If a selection is made for a locked game, it will be ignored.
func postPicks(user nflpickem.User, db pickManager, notifier nflpickem.Notifier, w http.ResponseWriter, r *http.Request) {
	yearStr := r.FormValue("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "year query parameter must be integer")
		return
	}

	weekStr := r.FormValue("week")
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "week query parameter must be integer")
		return
	}

	username := r.FormValue("username")
	if username == "" {
		WriteJSONError(w, http.StatusBadRequest, "username is required")
		return
	}

	picks, err := db.UserPicks(username, year, week)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	selections := make(nflpickem.PickSet, 0)
	err = json.NewDecoder(r.Body).Decode(&selections)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Drop any games that have already started
	selections = selections.Filter(func(p nflpickem.Pick) bool {
		return p.Game.Date.After(time.Now())
	})

	err = picks.Merge(selections)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if !picks.IsLegal() {
		WriteJSONError(w, http.StatusBadRequest, "resulting pick set contains too many point values")
		return
	}

	err = db.MakePicks(picks)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = notifier.Notify(username, week, picks)
	if err != nil {
		log.Printf("unable to notify user of picks: %v", err)
		return
	}

	WriteJSON(w, picks)
}
