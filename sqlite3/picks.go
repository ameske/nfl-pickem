package sqlite3

import (
	"errors"
	"log"
	"time"

	"github.com/ameske/nfl-pickem"
)

// ErrGameLocked occurs when a pick's game has already kicked off and cannot be changed
var ErrGameLocked = errors.New("game is locked")

// SelectedPicks returns the user's selected picks for the given week of the requested NFL season.
func (db Datastore) SelectedPicks(username string, year int, week int) (nflpickem.PickSet, error) {
	sql := `SELECT years.year, weeks.week, home.city, home.nickname, away.city, away.nickname, games.date, games.home_score, games.away_score, selection.city, selection.nickname, picks.points, users.first_name, users.last_name, users.email
		FROM picks
		JOIN games ON picks.game_id = games.id
		JOIN teams AS home ON games.home_id = home.id
		JOIN teams AS away ON games.away_id = away.id
		JOIN teams AS selection ON picks.selection = selection.id
		JOIN weeks ON games.week_id = weeks.id
		JOIN years ON weeks.year_id = years.id
		JOIN users ON picks.user_id = users.id
		WHERE picks.selection IS NOT NULL AND users.email LIKE ?1 AND years.year = ?2 AND weeks.week = ?3`

	rows, err := db.Query(sql, username, year, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	picks := make(nflpickem.PickSet, 0)

	for rows.Next() {
		var tmp nflpickem.Pick
		var d int64
		err := rows.Scan(&tmp.Game.Year, &tmp.Game.Week, &tmp.Game.Home.City, &tmp.Game.Home.Nickname, &tmp.Game.Away.City, &tmp.Game.Away.Nickname, &d, &tmp.Game.HomeScore, &tmp.Game.AwayScore,
			&tmp.Selection.City, &tmp.Selection.Nickname,
			&tmp.Points,
			&tmp.User.FirstName, &tmp.User.LastName, &tmp.User.Email)
		if err != nil {
			return nil, err
		}

		tmp.Game.Date = time.Unix(d, 0)

		picks = append(picks, tmp)
	}

	return picks, nil
}

// SelectedPicks returns the user's selected picks for the given week of the requested NFL season.
func (db Datastore) UnselectedPicks(username string, year int, week int) (nflpickem.PickSet, error) {
	sql := `SELECT years.year, weeks.week, home.city, home.nickname, away.city, away.nickname, games.date, games.home_score, games.away_score, users.first_name, users.last_name, users.email
		FROM picks
		JOIN games ON picks.game_id = games.id
		JOIN teams AS home ON games.home_id = home.id
		JOIN teams AS away ON games.away_id = away.id
		JOIN weeks ON games.week_id = weeks.id
		JOIN years ON weeks.year_id = years.id
		JOIN users ON picks.user_id = users.id
		WHERE picks.selection IS NULL AND users.email LIKE ?1 AND years.year = ?2 AND weeks.week = ?3`

	rows, err := db.Query(sql, username, year, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	picks := make(nflpickem.PickSet, 0)

	for rows.Next() {
		var tmp nflpickem.Pick
		var d int64
		err := rows.Scan(&tmp.Game.Year, &tmp.Game.Week, &tmp.Game.Home.City, &tmp.Game.Home.Nickname, &tmp.Game.Away.City, &tmp.Game.Away.Nickname, &d, &tmp.Game.HomeScore, &tmp.Game.AwayScore,
			&tmp.User.FirstName, &tmp.User.LastName, &tmp.User.Email)
		if err != nil {
			return nil, err
		}

		tmp.Game.Date = time.Unix(d, 0)

		picks = append(picks, tmp)
	}

	return picks, nil
}

// Picks returns all picks for a given week of the requested NFL season
func (db Datastore) Picks(year int, week int) (nflpickem.PickSet, error) {
	selected, err := db.SelectedPicks("%", year, week)
	if err != nil {
		return nil, err
	}

	unselected, err := db.UnselectedPicks("%", year, week)
	if err != nil {
		return nil, err
	}

	return append(selected, unselected...), nil
}

// Picks returns the given user's picks for the given week of the requested NFL season.
func (db Datastore) UserPicks(username string, year int, week int) (nflpickem.PickSet, error) {
	selected, err := db.SelectedPicks("%", year, week)
	if err != nil {
		return nil, err
	}

	unselected, err := db.UnselectedPicks("%", year, week)
	if err != nil {
		return nil, err
	}

	return append(selected, unselected...), nil
}

// TODO: Implement MakePicks
func (db Datastore) MakePicks(picks nflpickem.PickSet) error {
	log.Println(picks)
	return nil
}

func (db Datastore) CreatePicks(username string, year int, week int) error {
	games, err := gameIds(db, year, week)
	if err != nil {
		return nil
	}

	sql := `INSERT INTO picks(user_id, game_id) VALUES((SELECT id FROM users WHERE email = ?1), ?2)`

	for _, gid := range games {
		_, err = db.Exec(sql, username, gid)
		if err != nil {
			return err
		}
	}

	return nil
}

func updatePick(db Datastore, selection int, points int, id int) error {
	sql := `UPDATE picks
	  SET selection = ?1 AND points = ?2
	  WHERE id = ?3`

	_, err := db.Exec(sql, selection, points, id)

	return err
}

func gameIds(db Datastore, year int, week int) ([]int, error) {
	sql := `SELECT id
		FROM games
		WHERE week_id = (SELECT weeks.id
				  FROM weeks
				  JOIN years ON weeks.year_id = years.id
				  WHERE years.year = ?1 AND weeks.week = ?2)`

	rows, err := db.Query(sql, year, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := make([]int, 0)

	for rows.Next() {
		var tmp int
		err = rows.Scan(&tmp)
		if err != nil {
			return nil, err
		}

		games = append(games, tmp)
	}

	return games, rows.Err()
}
