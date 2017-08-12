package sqlite3

import (
	"time"

	"github.com/ameske/nfl-pickem"
)

// Results returns the set of picks that are visible for all users for the given week of the NFL season.
//
// A pick is visible if the game is locked.
func (db Datastore) Results(t time.Time, year int, week int) ([]nflpickem.Result, error) {
	sql := `SELECT years.year, weeks.week, home.city, home.nickname, away.city, away.nickname, games.date, games.home_score, games.away_score, selection.city, selection.nickname, picks.points, users.first_name, users.last_name, users.email
		FROM picks
		JOIN games ON picks.game_id = games.id
		JOIN teams AS home ON games.home_id = home.id
		JOIN teams AS away ON games.away_id = away.id
		JOIN teams AS selection ON picks.selection = selection.id
		JOIN weeks ON games.week_id = weeks.id
		JOIN years ON weeks.year_id = years.id
		JOIN users ON picks.user_id = users.id
		WHERE picks.selection IS NOT NULL AND games.date < ?1 AND years.year = ?2 AND weeks.week = ?3 ORDER BY games.date ASC, games.id ASC, users.email ASC`

	rows, err := db.Query(sql, t.Unix(), year, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seenGames := make(map[nflpickem.Game]bool)
	results := make([]nflpickem.Result, 0)
	current := -1

	for rows.Next() {
		var g nflpickem.Game
		var pr nflpickem.PickResult
		var d int64

		err := rows.Scan(&g.Year, &g.Week, &g.Home.City, &g.Home.Nickname, &g.Away.City, &g.Away.Nickname, &d, &g.HomeScore, &g.AwayScore, &pr.Selection.City, &pr.Selection.Nickname, &pr.Points, &pr.User.FirstName, &pr.User.LastName, &pr.User.Email)
		if err != nil {
			return nil, err
		}

		g.Date = time.Unix(d, 0)

		if !seenGames[g] {
			seenGames[g] = true
			current++
			results = append(results, nflpickem.Result{Game: g, Picks: make([]nflpickem.PickResult, 0)})
		}

		results[current].Picks = append(results[current].Picks, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
