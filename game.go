package nflpickem

import "time"

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

func (g Game) Equal(other Game) bool {
	return (g.Year == other.Year && g.Week == other.Week &&
		g.Date.Equal(other.Date) && g.Home.Equal(other.Home) &&
		g.Away.Equal(other.Away))
}

// GamesRetriever is the interface implemented by a type that can retrieve NFL game
// information.
type GamesRetriever interface {
	WeekGames(year int, week int) ([]Game, error)
	CumulativeGames(year int, week int) ([]Game, error)
}

type Updater interface {
	Weeker
	UpdateGame(week int, year int, homeTeam string, homeScore int, awayScore int) error
}

// Team represents an NFL team
type Team struct {
	City     string `json:"city"`
	Nickname string `json:"nickname"`
}

func (t Team) Equal(other Team) bool {
	return t.City == other.City && t.Nickname == other.Nickname
}
