package sqlite3

import "database/sql"

func (db Datastore) teamRecord(city string, nickname string) (wins int, losses int, err error) {
	// If somebody doesn't have any wins or losses, we have trouble coming up with their record all at once.
	// For now, until I figure out what I'm doing wrong, calculate it individually.
	var homeWins, homeLosses, awayWins, awayLosses int

	s_home_wins := `SELECT COUNT(*) FROM games JOIN teams ON games.home_id = teams.id WHERE home_score > away_score AND teams.city = ?1 AND teams.nickname = ?2`
	err = db.QueryRow(s_home_wins, city, nickname).Scan(&homeWins)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	} else if err != nil {
		return -1, -1, err
	}

	s_away_wins := `SELECT COUNT(*) FROM games JOIN teams ON games.away_id = teams.id WHERE away_score > home_score AND teams.city = ?1 AND teams.nickname = ?2`
	err = db.QueryRow(s_away_wins, city, nickname).Scan(&awayWins)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	} else if err != nil {
		return -1, -1, err
	}

	s_home_losses := `SELECT COUNT(*) FROM games JOIN teams ON games.home_id = teams.id WHERE home_score < away_score AND teams.city = ?1 AND teams.nickname = ?2`
	err = db.QueryRow(s_home_losses, city, nickname).Scan(&homeLosses)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	} else if err != nil {
		return -1, -1, err
	}

	s_away_losses := `SELECT COUNT(*) FROM games JOIN teams ON games.away_id = teams.id WHERE away_score < home_score AND teams.city = ?1 AND teams.nickname = ?2`
	err = db.QueryRow(s_away_losses, city, nickname).Scan(&awayLosses)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	} else if err != nil {
		return -1, -1, err
	}

	return homeWins + awayWins, homeLosses + awayLosses, nil
}
