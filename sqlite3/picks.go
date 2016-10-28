package sqlite3

import (
	"database/sql"
	"errors"
	"time"

	"github.com/ameske/nfl-pickem"
)

const (
	PICK_CORRECT   = 1
	PICK_INCORRECT = 0
)

var ErrGameLocked = errors.New("game is locked")

func processPickResults(rows *sql.Rows) (nflpickem.PickSet, error) {
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

func (db Datastore) SelectedPicks(username string, year int, week int) (nflpickem.PickSet, error) {
	sql := `SELECT years.year, weeks.week, home.city, home.nickname, away.city, away.nickname, games.date, games.home_score, games.away_score, selection.city, selection.nickname, picks.points, users.first_name, users.last_name, users.email
		FROM picks
		JOIN games ON picks.game_id = games.id
		JOIN teams AS home ON games.home_id = home.id
		JOIN teams AS away ON games.away_id = away.id
		JOIN teams AS selection ON picks.selection_id = selection.id
		JOIN weeks ON games.week_id = weeks.id
		JOIN years ON weeks.year_id = years.id
		JOIN users ON picks.user_id = users.id
		WHERE picks.selection IS NOT NULL AND users.email = ?1 AND years.year = ?2 AND weeks.week = ?3`

	rows, err := db.Query(sql, username, year, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return processPickResults(rows)
}

func (db Datastore) Picks(username string, year int, week int) (nflpickem.PickSet, error) {
	sql := `SELECT years.year, weeks.week, home.city, home.nickname, away.city, away.nickname, games.date, games.home_score, games.away_score, selection.city, selection.nickname, picks.points, users.first_name, users.last_name, users.email
		FROM picks
		JOIN games ON picks.game_id = games.id
		JOIN teams AS home ON games.home_id = home.id
		JOIN teams AS away ON games.away_id = away.id
		JOIN teams AS selection ON picks.selection_id = selection.id
		JOIN weeks ON games.week_id = weeks.id
		JOIN years ON weeks.year_id = years.id
		JOIN users ON picks.user_id = users.id
		WHERE users.email = ?1 AND years.year = ?2 AND weeks.week = ?3`

	rows, err := db.Query(sql, username, year, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return processPickResults(rows)
}

func (db Datastore) Grade(id int, points int, correct bool) error {
	var intBool int
	if correct {
		intBool = PICK_CORRECT
	} else {
		intBool = PICK_INCORRECT
	}

	_, err := db.Exec("UPDATE picks SET correct = ?1, points = ?2 WHERE id = ?3", intBool, points, id)

	return err
}

func (db Datastore) MakePicks(picks nflpickem.PickSet) error {
	return nil
}
