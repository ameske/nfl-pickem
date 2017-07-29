package sqlite3

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Datastore is a handle to a sqlite3 database storing NFL pickem data.
type Datastore struct {
	*sql.DB
}

// NewDatastore connects to a sqlite3 database storing NFL pickem data.
func NewDatastore(path string) (*Datastore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Datastore{DB: db}, nil
}
