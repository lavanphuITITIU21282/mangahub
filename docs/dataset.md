# Manga Dataset

This project ships with a small starter dataset so the API and CLI have enough titles to demonstrate search, details, library, and progress features.

## Files

- `data/manga_seed.json` — human-friendly dataset (48 manga) that can be extended.
- `internal/db/seed_manga.sql` — generated seed script used by migrations.

## How it loads

`internal/db/migrate.go` embeds `schema.sql` and `seed_manga.sql`. Every time the server starts and calls `db.Migrate(...)`, the seed runs with `INSERT OR IGNORE`, so it is safe to run multiple times.

## Extending the dataset

1. Add new entries to `data/manga_seed.json`.
2. Regenerate `internal/db/seed_manga.sql` (optional). In this repo, the seed SQL is kept in sync with the JSON.

