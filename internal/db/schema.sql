CREATE TABLE IF NOT EXISTS manga (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  author TEXT NOT NULL,
  genres TEXT NOT NULL,
  status TEXT NOT NULL,
  total_chapters INTEGER NOT NULL,
  description TEXT NOT NULL
);
