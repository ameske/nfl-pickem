package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ameske/nfl-pickem"
)

// Games returns the JSON representation of NFL games.
//
// URL Parameters:
//	year: Specifies the current year, Required
//	week: Specifies the current week, Required
//	kind: ["cumulative"], returns games for the current year up to the given week, Optional
func games(db nflpickem.GamesRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var games []nflpickem.Game

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

		kind := r.FormValue("kind")

		switch kind {
		case "":
			games, err = db.WeekGames(year, week)
		case "cumulative":
			games, err = db.CumulativeGames(year, week)
		default:
			WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("unknown kind parameter value [%s]", kind))
			return
		}

		if err != nil {
			WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		WriteJSON(w, games)
	}
}
