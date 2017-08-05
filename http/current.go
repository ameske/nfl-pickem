package http

import (
	"net/http"
	"time"

	"github.com/ameske/nfl-pickem"
)

// currentWeek returns the JSON representation of the current week of the NFL season
func currentWeek(db nflpickem.Weeker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		week, err := db.CurrentWeek(time.Now())
		if err != nil {
			WriteJSONError(w, http.StatusInternalServerError, err.Error())
		}

		WriteJSON(w, week)
	}
}

func years(db nflpickem.DataSummarizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		years, err := db.Years()
		if err != nil {
			WriteJSONError(w, http.StatusInternalServerError, err.Error())
		}

		y := struct {
			Years []int `json:"years"`
		}{
			years,
		}
		WriteJSON(w, y)
	}
}
