package sqlite3

import (
	"database/sql"
	"time"

	"github.com/ameske/nfl-pickem"
)

const (
	oneWeek      = time.Hour * 24 * 7
	seasonLength = 17
)

// CurrentWeek returns the current week of the season. A season starts on the
// Tuesday before the first game. A new week starts every Tuesday. Given the
// current time, we can calculate the current week of the season. A week set
// to -1 means that we are in the offseason.
func (db Datastore) CurrentWeek(t time.Time) (nflpickem.Week, error) {
	start, err := db.currentSeasonStart(t)
	if err != nil {
		return nflpickem.Week{Year: -1, Week: -1}, err
	}

	d := t.Sub(start)

	week := int(d/oneWeek) + 1

	if week > seasonLength {
		return nflpickem.Week{Year: start.Year(), Week: -1}, nil
	}

	return nflpickem.Week{Year: start.Year(), Week: week}, nil
}

func (db Datastore) currentSeasonStart(t time.Time) (start time.Time, err error) {
	now := t.Unix()

	var s sql.NullInt64
	row := db.QueryRow("SELECT MAX(year_start) FROM years WHERE year_start < ?1", now)
	err = row.Scan(&s)
	if err != nil {
		return time.Unix(0, 0), err
	}

	// Special case: if now + 7 is a different value then that means we're on the cusp of a new season. So pretend we are in week 1.
	now2 := time.Date(t.Year(), t.Month(), t.Day()+7, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	var s2 sql.NullInt64
	row = db.QueryRow("SELECT MAX(year_start) FROM years WHERE year_start < ?1", now2.Unix())
	err = row.Scan(&s2)
	if err != nil {
		return time.Unix(0, 0), err
	}

	if s.Int64 != s2.Int64 && s2.Int64 != 0 {
		return time.Unix(s2.Int64-604800, 0), err
	} else if s.Valid {
		return time.Unix(s.Int64, 0), err
	} else {
		return time.Unix(0, 0), err
	}
}

// AddWeek adds the week, associated with the given year to the datastore.
func (db Datastore) AddWeek(year int, week int) error {
	_, err := db.Exec("INSERT INTO weeks(week, year_id) VALUES(?1, (SELECT id FROM YEARS where year = ?2))", week, year)

	return err
}

// AddYear adds the year with the given start epoch to the datastore.
func (db Datastore) AddYear(year int, yearStart int) error {
	_, err := db.Exec("INSERT INTO years(year, year_start) VALUES(?1, ?2)", year, yearStart)

	return err
}
