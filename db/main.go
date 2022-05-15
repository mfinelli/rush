package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func SetupDB() error {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		return err
	}

	rdb := NewRushDB(db)

	if err := rdb.Migrate(); err != nil {
		return err
	}

	return nil
}
