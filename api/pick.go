package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ameske/nfl-pickem/jsonhttp"
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

// PickRetriever is the interface implemented by types that can retrieve Pick information
type PickRetriever interface {
	Picks(username string, year int, week int) ([]Pick, error)
	SelectedPicks(username string, year int, week int) ([]Pick, error)
}

// GetPicks returns the set of picks for the given user, year, and week.
//
// Using the "kind" parameters it is possible to specify that all picks
// should be returned, or only selected picks.
func GetPicks(db PickRetriever) http.HandlerFunc {
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
// The selection field is the nickname of the given team. This is enough to uniquely identify the game.
//
// The set of (username, year, week, selection) is enough to uniquely identify the backing pick if we are allowing
// more than one user to submit picks at a time (a la an admin form of some sort).
type Selection struct {
	Username  string `json:"username"`
	Year      int    `json:"year"`
	Week      int    `json:"week"`
	Selection Team   `json:"selection"`
	Points    int    `json:"points"`
}

var (
	errGameLocked       = errors.New("game has already started - pick locked")
	errUnknownSelection = errors.New("selection does not match a game in the given pick set")
)

// Picker is the interface implemented by a type that can make/update picks
type Picker interface {
	PickRetriever
	MakePicks([]Pick) error
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

		picks, err := db.Picks(username, year, week)
		if err != nil {
			jsonhttp.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Manipulate a represenation in memory first before we commit to the backing datastore
		now := time.Now()
		for i, s := range selections {
			if s.Username != username {
				jsonhttp.WriteError(w, http.StatusBadRequest, "selections must match the username declared in the URL parameter")
				return
			}

			err := updateTemporaryPicks(now, s, picks)
			if err == errGameLocked {
				continue
			} else if err == errUnknownSelection {
				jsonhttp.WriteError(w, http.StatusForbidden, fmt.Sprintf("selection %d is invalid for Year %d / Week %d - unable to update", i, year, week))
				return
			} else {
				jsonhttp.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		if !pickSetLegal(picks) {
			jsonhttp.WriteError(w, http.StatusBadRequest, "resulting pick set contains too many point values")
			return
		}

		err = db.MakePicks(picks)
		if err != nil {
			jsonhttp.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		jsonhttp.Write(w, picks)
	}
}

// updateTemporaryPicks locates the correct pick given the selection and changes the selection and points
// provided that the game time is not before the provided time.
func updateTemporaryPicks(t time.Time, s Selection, picks []Pick) error {
	for i, p := range picks {
		// We need to match on (year, week, home|away team) to know we can update the pick with the given selection
		if s.Year == p.Game.Year && s.Week == p.Game.Week && (s.Selection.Nickname == p.Game.Home.Nickname || s.Selection.Nickname == p.Game.Away.Nickname) {
			// Refuse to update the pick if the game has already started
			if p.Game.Date.Before(t) {
				return errGameLocked
			}
			picks[i].Points = s.Points
			picks[i].Selection.City = s.Selection.City
			picks[i].Selection.Nickname = s.Selection.Nickname
			return nil
		}
	}

	return errUnknownSelection
}

const (
	maxSevens = 1
	maxFives  = 2
	maxThrees = 5
)

// pickSetLegal determines the legality of a set of picks, assuming that the set of
// picks represents only one (year, week, user) set.
func pickSetLegal(picks []Pick) bool {
	threes := 0
	fives := 0
	sevens := 0

	for _, p := range picks {
		switch p.Points {
		case 3:
			threes++
		case 5:
			fives++
		case 7:
			sevens++

		}

	}

	return threes <= maxThrees && fives <= maxFives && sevens <= maxSevens
}
