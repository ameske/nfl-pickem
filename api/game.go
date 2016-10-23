package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ameske/nfl-pickem/jsonhttp"
)

// Game represents an NFL contest.
type Game struct {
	Year      int       `json:"year"`
	Week      int       `json:"week"`
	Date      time.Time `json:"date"`
	Home      Team      `json:"home"`
	Away      Team      `json:"away"`
	HomeScore int       `json:"homeScore"`
	AwayScore int       `json:"awayScore"`
}

// GamesRetriever is the interface implemented by a type that can retrieve NFL game
// information.
type GamesRetriever interface {
	WeekGames(year int, week int) ([]Game, error)
	CumulativeGames(year int, week int) ([]Game, error)
}

// Games returns the JSON representation of NFL games.
//
// URL Parameters:
//	year: Specifies the current year, Required
//	week: Specifies the current week, Required
//	kind: ["cumulative"], returns games for the current year up to the given week, Optional
func Games(db GamesRetriever) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var games []Game

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

		kind := r.FormValue("kind")

		switch kind {
		case "":
			games, err = db.WeekGames(year, week)
		case "cumulative":
			games, err = db.CumulativeGames(year, week)
		default:
			jsonhttp.WriteError(w, http.StatusBadRequest, fmt.Sprintf("unknown kind parameter value [%s]", kind))
			return
		}

		if err != nil {
			jsonhttp.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		jsonhttp.Write(w, games)
	}
}
