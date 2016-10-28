package nflpickem

import "time"

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
