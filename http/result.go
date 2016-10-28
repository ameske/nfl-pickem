package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ameske/nfl-pickem"
)

// Results returns the set of picks for the given week, where the game has already started.
//
// This endpoint sorts the games by date, and sorts the list of pick results by username.
func results(db nflpickem.ResultFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		results, err := db.Results(time.Now(), year, week)
		if err != nil {
			WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		WriteJSON(w, results)
	}
}
