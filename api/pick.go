package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ameske/nfl-pickem/jsonhttp"
)

const (
	SELECTION_AWAY = 1
	SELECTION_HOME = 2
)

// A Pick represents a user's selection for a given game.
//
// The Pick can stand on its own since it contains embedded game information
type Pick struct {
	Game      Game `json:"game"`
	User      User `json:"user"`
	Selection Team `json:"selection"`
	Points    int  `json:"points"`
}

// Picker is the interface implemented by types that can retrieve Pick information
type Picker interface {
	Picks(username string, year int, week int) ([]Pick, error)
	SelectedPicks(username string, year int, week int) ([]Pick, error)
}

// GetPicks returns the set of picks for the given user, year, and week.
//
// Using the "kind" parameters it is possible to specify that all picks
// should be returned, or only selected picks.
func GetPicks(db Picker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		yearStr := r.FormValue("year")
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			jsonhttp.WriteError(w, http.StatusBadRequest, "year query parameter must be integer")
			return
		}

		weekStr := r.FormValue("week")
		week, err := strconv.Atoi(weekStr)
		if err != nil {
			jsonhttp.WriteError(w, http.StatusBadRequest, "week query parameter must be integer")
			return
		}

		username := r.FormValue("username")
		kind := r.FormValue("kind")

		var picks []Pick
		switch kind {
		case "":
			picks, err = db.Picks(username, year, week)
		case "all":
			picks, err = db.Picks(username, year, week)
		case "selected":
			picks, err = db.SelectedPicks(username, year, week)
		default:
			jsonhttp.WriteError(w, http.StatusBadRequest, fmt.Sprintf("unknown kind value: %s", kind))
			return
		}

		if err != nil {
			jsonhttp.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		jsonhttp.Write(w, picks)
	}
}

// A Selection represents a completed pick by a user.
//
// The set of (year, week, selection) is enough to uniquely identify the backing pick.
type Selection struct {
	Username  string `json:"username"`
	Year      int    `json:"year"`
	Week      int    `json:"week"`
	Selection string `json:"selection"`
	Points    int    `json:"points"`
}

const (
	maxSeven = 1
	maxFive  = 2
	maxThree = 5
)

// MakePicks processes an array of JSON representation of pick selections.
//
// In the event duplicate picks for the same game are made,
// the last pick is always the pick that is stored.
//
// Any picks that are locked will be ignored by the endpoint.
func MakePicks(db Picker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		yearStr := r.FormValue("year")
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			jsonhttp.WriteError(w, http.StatusBadRequest, "year query parameter must be integer")
			return
		}

		weekStr := r.FormValue("week")
		week, err := strconv.Atoi(weekStr)
		if err != nil {
			jsonhttp.WriteError(w, http.StatusBadRequest, "week query parameter must be integer")
			return
		}

		// For now we'll restrict this to one user, the selection type enables us to have multiple
		// user picks made for admin functionality one day.
		username := r.FormValue("username")
		if username == "" {
			jsonhttp.WriteError(w, http.StatusBadRequest, "username is required")
			return
		}

		selections := make([]Selection, 0)
		err = json.NewDecoder(r.Body).Decode(&selections)
		if err != nil {
			jsonhttp.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Gather the current set of picks. We'll store this in memory and manipulate
		// this to determine if we can commit the transaction.
		_, err = db.Picks(username, year, week)
		if err != nil {
			jsonhttp.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

// updatePick locates the correct pick given the selection and changes the selection and points
func updatePick(s Selection, picks []Pick) {
}
