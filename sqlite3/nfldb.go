package sqlite3

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Datastore struct {
	*sql.DB
}

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
