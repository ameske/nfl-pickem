package sqlite3

import (
	"time"

	"github.com/ameske/nfl-pickem"
)

// WeekGames returns the games for only the specified week from the datastore.
func (db Datastore) WeekGames(year int, week int) ([]nflpickem.Game, error) {
	return db.games(year, week, week)
}

// CumulativeGames returns games up to the specified week from the datastore.
func (db Datastore) CumulativeGames(year int, week int) ([]nflpickem.Game, error) {
	return db.games(year, 1, week)
}

func (db Datastore) games(year int, minWeek int, maxWeek int) ([]nflpickem.Game, error) {
	sql := `SELECT years.year, weeks.week, games.date, home.city, home.nickname, away.city, away.nickname, games.home_score, games.away_score
	    FROM games
	    JOIN teams AS home ON games.home_id = home.id
	    JOIN teams AS away ON games.away_id = away.id
	    JOIN weeks ON games.week_id = weeks.id
	    JOIN years ON weeks.year_id = years.id
	    WHERE years.year = ?1 AND weeks.week >= ?2 AND weeks.week <= ?3`

	rows, err := db.Query(sql, year, minWeek, maxWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := make([]nflpickem.Game, 0)

	for rows.Next() {
		var tmp nflpickem.Game
		var d int64

		err := rows.Scan(&tmp.Year, &tmp.Week, &d, &tmp.Home.City, &tmp.Home.Nickname, &tmp.Away.City, &tmp.Away.Nickname, &tmp.HomeScore, &tmp.AwayScore)
		if err != nil {
			return nil, err
		}

		tmp.Date = time.Unix(d, 0)

		games = append(games, tmp)
	}

	return games, nil
}

// UpdateGame overwrites the stored home and away score for the given game. The game is specified by the
// year, week, and home team, which uniquely identifies any game.
func (db Datastore) UpdateGame(week int, year int, homeTeam string, homeScore int, awayScore int) error {
	// sqlite3 makes this hard on us by not allowing JOIN in UPDATE
	// so we have to do this in a couple of steps
	sql := `SELECT games.id FROM games
		JOIN weeks ON games.week_id = weeks.id
		JOIN years ON weeks.year_id = years.id
		JOIN teams ON games.home_id = teams.id
		WHERE weeks.week = ?1 AND years.year = ?2 AND teams.nickname = ?3`

	var gameId int64
	err := db.QueryRow(sql, week, year, homeTeam).Scan(&gameId)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE games
			  SET home_score = ?2, away_score = ?3
			  WHERE id = ?1`, gameId, homeScore, awayScore)

	return err
}

// AddGame adds the given game to the datastore.
func (db Datastore) AddGame(date time.Time, homeTeam string, awayTeam string) error {
	nflWeek, err := db.CurrentWeek(date)
	if err != nil {
		return err
	}

	// If the month is January, the NFL regular season bled over into the next calendar year, we need to adjust the year
	// back one so that the game is counted as part of the correct NFL year.
	if date.Month() == 1 {
		_, err = db.Exec(`INSERT INTO games(week_id, date, home_id, away_id)
			 VALUES((SELECT weeks.id FROM weeks JOIN years ON weeks.year_id = years.id WHERE years.year = ?1 AND weeks.week = ?2), ?3, (SELECT id FROM teams WHERE nickname = ?4), (SELECT id FROM teams WHERE nickname = ?5))`, date.Year()-1, nflWeek.Week, date.Unix(), homeTeam, awayTeam)
	} else {
		_, err = db.Exec(`INSERT INTO games(week_id, date, home_id, away_id)
			 VALUES((SELECT weeks.id FROM weeks JOIN years ON weeks.year_id = years.id WHERE years.year = ?1 AND weeks.week = ?2), ?3, (SELECT id FROM teams WHERE nickname = ?4), (SELECT id FROM teams WHERE nickname = ?5))`, date.Year(), nflWeek.Week, date.Unix(), homeTeam, awayTeam)

	}

	return err
}
