package db

import (
	"database/sql"
	_ "embed"
)

//go:embed schema.sql
var schema string

//go:embed seed_manga.sql
var seedManga string

func Migrate(db *sql.DB) error {
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	// Seed manga dataset (safe to run multiple times).
	if _, err := db.Exec(seedManga); err != nil {
		return err
	}

	return nil
}
