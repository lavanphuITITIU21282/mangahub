CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS manga (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  author TEXT NOT NULL,
  genres TEXT NOT NULL,
  status TEXT NOT NULL,
  total_chapters INTEGER NOT NULL,
  description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_progress (
  user_id TEXT NOT NULL,
  manga_id TEXT NOT NULL,
  current_chapter INTEGER NOT NULL,
  status TEXT NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, manga_id)
);
