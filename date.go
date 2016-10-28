package nflpickem

import "time"

// Week represents a unique week of the NFL Pickem' Pool
type Week struct {
	Year int `json:"year"`
	Week int `json:"week"`
}

// Week represents the current week of an NFL season.
// Weeker is the interface implemented by types who can retrieve the current week
// of the NFL season given a point in time.
type Weeker interface {
	CurrentWeek(t time.Time) (w Week, err error)
}
