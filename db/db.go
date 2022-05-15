package db

import (
	"database/sql"
	// _ "github.com/mattn/go-sqlite3"
)

type RushDB struct {
	db *sql.DB
}

func NewRushDB(db *sql.DB) *RushDB {
	return &RushDB{db: db}
}

func (r *RushDB) Migrate() error {
	query := `CREATE TABLE IF NOT EXISTS ca_certificates(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		public_key TEXT NOT NULL,
		private_Key TEXT NOT NULL);`

	_, err := r.db.Exec(query)
	return err
}
