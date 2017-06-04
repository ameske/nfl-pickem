package sqlite3

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (db Datastore) IsAdmin(username string) (admin bool, err error) {
	row := db.QueryRow("SELECT admin FROM users WHERE email = ?1", username)
	err = row.Scan(&admin)
	if err != nil {
		return false, err
	}

	return
}

func (db Datastore) UserFirstNames() ([]string, error) {
	rows, err := db.Query("SELECT first_name FROM users ORDER BY id ASC")
	if err != nil {
		return nil, err
	}

	users := make([]string, 0)
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			return nil, err
		}
		users = append(users, tmp)
	}
	rows.Close()

	return users, nil
}

func (db Datastore) Usernames() ([]string, error) {
	rows, err := db.Query("SELECT email FROM users ORDER BY id ASC")
	if err != nil {
		return nil, err
	}

	users := make([]string, 0)
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			return nil, err
		}
		users = append(users, tmp)
	}
	rows.Close()

	return users, nil
}

func (db Datastore) CheckCredentials(user string, password string) (bool, error) {
	var storedPassword string
	row := db.QueryRow("SELECT password FROM users WHERE email = ?1", user)
	err := row.Scan(&storedPassword)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))

	return err == nil, nil
}

func (db Datastore) UpdatePassword(user string, oldPassword string, newPassword string) error {
	ok, err := db.CheckCredentials(user, oldPassword)
	if err != nil {
		return err
	} else if !ok {
		return errors.New("unauthorized")
	}

	_, err = db.Exec("UPDATE users SET password = ?1 WHERE email = ?2", string(newPassword), user)
	return err
}
