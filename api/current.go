package api

import (
	"net/http"
	"time"

	"github.com/ameske/nfl-pickem/jsonhttp"
)

// Week represents the current week of an NFL season.
type Week struct {
	Year int `json:"year"`
	Week int `json:"week"`
}

// Weeker is the interface implemented by types who can retrieve the current week
// of the NFL season given a point in time.
type Weeker interface {
	CurrentWeek(t time.Time) (year int, week int, err error)
}

// CurrentWeek writes the JSON representation of the current week of the NFL season
func CurrentWeek(db Weeker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		year, week, err := db.CurrentWeek(time.Now())
		if err != nil {
			jsonhttp.WriteError(w, http.StatusInternalServerError, err.Error())
		}

		jsonhttp.Write(w, Week{Year: year, Week: week})
	}
}
