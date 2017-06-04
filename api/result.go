package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ameske/nfl-pickem/jsonhttp"
)

// A Result consists of a game, and all of the picks made by users for that game.
type Result struct {
	Game  Game         `json:"game"`
	Picks []PickResult `json:"picks"`
}

// PickResult contains just the pick information for a given game and user.
//
// A PickResult is not very useful on its own, and is always returned as part of a Result so
// that the client is able to associate it with a game.
type PickResult struct {
	User      User `json:"user"`
	Selection Team `json:"selection"`
	Points    int  `json:"points"`
}

// ResultFetcher is the interface implemented by types that can fetch results for a given year,
// and week of a season.
type ResultFetcher interface {
	Results(t time.Time, year int, week int) ([]Result, error)
}

// Results returns the set of picks for the given week, where the game has already started.
//
// This endpoint sorts the games by date, and sorts the list of pick results by username.
func Results(db ResultFetcher) http.HandlerFunc {
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

		results, err := db.Results(time.Now(), year, week)
		if err != nil {
			jsonhttp.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		jsonhttp.Write(w, results)
	}
}
