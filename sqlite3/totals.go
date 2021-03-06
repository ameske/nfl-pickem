package sqlite3

import "github.com/ameske/nfl-pickem"

// UserWeekTotal returns the user's total for the given week of the NFL season.
func (db Datastore) UserWeekTotal(username string, year int, week int) ([]nflpickem.WeekTotal, error) {
	return db.weekTotals(username, year, week, week)
}

// UserWeekTotals returns the totals for all weeks up to the given week of the NFL season.
func (db Datastore) UserWeekTotals(username string, year int, week int) ([]nflpickem.WeekTotal, error) {
	return db.weekTotals(username, year, 1, week)
}

// WeekTotals reutrns all users totals for the given week of the NFL season.
func (db Datastore) WeekTotals(year int, week int) ([]nflpickem.WeekTotal, error) {
	return db.weekTotals("%", year, week, week)
}

// CumulativeWeekTotals returns all users totals up to the given week of the NFL season.
func (db Datastore) CumulativeWeekTotals(year int, week int) ([]nflpickem.WeekTotal, error) {
	return db.weekTotals("%", year, 1, week)
}

func (db Datastore) weekTotals(username string, year int, minWeek int, maxWeek int) ([]nflpickem.WeekTotal, error) {
	sql := `SELECT users.first_name, users.last_name, users.email, years.year, weeks.week, SUM(picks.points)
		FROM picks
		JOIN users ON picks.user_id = users.id
		JOIN games ON picks.game_id = games.id
		JOIN weeks ON games.week_id = weeks.id JOIN years ON weeks.year_id = years.id
		WHERE users.email LIKE ?1 AND years.year = ?2 AND weeks.week >= ?3 AND weeks.week <= ?4 AND ((games.home_score > games.away_score AND picks.selection = games.home_id) OR (games.home_score < games.away_score AND picks.selection = games.away_id))
		GROUP BY users.email, weeks.week`

	rows, err := db.Query(sql, username, year, minWeek, maxWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totals := make([]nflpickem.WeekTotal, 0)

	for rows.Next() {
		var tmp nflpickem.WeekTotal
		err := rows.Scan(&tmp.User.FirstName, &tmp.User.LastName, &tmp.User.Email, &tmp.Year, &tmp.Week, &tmp.Total)
		if err != nil {
			return nil, err
		}

		totals = append(totals, tmp)
	}

	return totals, nil
}
