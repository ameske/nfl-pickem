package sqlite3

import (
	"github.com/ameske/nfl-pickem"

	"golang.org/x/crypto/bcrypt"
)

var unknownUser = nflpickem.User{}

func (db Datastore) CheckCredentials(username string, password string) (nflpickem.User, error) {
	var storedPassword string
	var user nflpickem.User

	row := db.QueryRow("SELECT users.first_name, users.last_name, users.email, users.admin, users.password FROM users WHERE email = ?1", username)
	err := row.Scan(&user.FirstName, &user.LastName, &user.Email, &user.Admin, &storedPassword)
	if err != nil {
		return unknownUser, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return unknownUser, err
	}

	return user, nil
}

func (db Datastore) UpdatePassword(username string, oldPassword string, newPassword string) error {
	_, err := db.CheckCredentials(username, oldPassword)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE users SET password = ?1 WHERE email = ?2", string(newPassword), username)
	return err
}
