package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ameske/nfl-pickem"
)

// WeeklyTotals returns the pre-calculated results for all users for a given year and week.
func weeklyTotals(db nflpickem.WeekTotalFetcher) http.HandlerFunc {
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

		kind := r.FormValue("type")

		var totals []nflpickem.WeekTotal

		switch kind {
		case "":
			totals, err = db.WeekTotals(year, week)
		case "cumulative":
			totals, err = db.CumulativeWeekTotals(year, week)
		default:
			WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("unknown kind parameter [%s]", kind))
			return
		}

		if err != nil {
			WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		WriteJSON(w, totals)
	}
}
