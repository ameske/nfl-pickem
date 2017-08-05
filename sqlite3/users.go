package sqlite3

import (
	"github.com/ameske/nfl-pickem"

	"golang.org/x/crypto/bcrypt"
)

var unknownUser = nflpickem.User{}

// CheckCredentials compares the given password to the store password hash in the datastore.
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

// UpdatePassword updates the given user's password in the datastore, hashing it before storing it.
func (db Datastore) UpdatePassword(username string, oldPassword string, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE users SET password = ?1 WHERE email = ?2", string(hash), username)
	return err
}

func (db Datastore) AddUser(first string, last string, email string, password string, admin bool) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users(first_name, last_name, email, password, admin) VALUES(?1, ?2, ?3, ?4, ?5)", first, last, email, hash, admin)

	return err
}
