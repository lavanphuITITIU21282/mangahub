package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func Migrate(db *sql.DB) error {
	// absolute path tới schema.sql
	b, err := os.ReadFile("C:/DC/mangahub/internal/db/schema.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(b))
	return err
}
