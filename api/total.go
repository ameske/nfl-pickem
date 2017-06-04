package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ameske/nfl-pickem/jsonhttp"
)

// WeekTotal is the aggregation of all correct picks for a user.
type WeekTotal struct {
	User  User `json:"user"`
	Year  int  `json:"year"`
	Week  int  `json:"week"`
	Total int  `json:"total"`
}

// WeekTotalFetcher is the interface implemented by types that can retrieve
// the aggregated results for a given week.
type WeekTotalFetcher interface {
	WeekTotals(year int, week int) ([]WeekTotal, error)
	CumulativeWeekTotals(year int, wee int) ([]WeekTotal, error)
}

// WeeklyTotals returns the pre-calculated results for all users for a given year and week.
func WeeklyTotals(db WeekTotalFetcher) http.HandlerFunc {
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

		kind := r.FormValue("type")

		var totals []WeekTotal

		switch kind {
		case "":
			totals, err = db.WeekTotals(year, week)
		case "cumulative":
			totals, err = db.CumulativeWeekTotals(year, week)
		default:
			jsonhttp.WriteError(w, http.StatusBadRequest, fmt.Sprintf("unknown kind parameter [%s]", kind))
			return
		}

		if err != nil {
			jsonhttp.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		jsonhttp.Write(w, totals)
	}
}
