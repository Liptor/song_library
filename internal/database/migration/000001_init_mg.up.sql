CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    releaseDate DATE,
    text TEXT,
    link TEXT,
    "group" TEXT NOT NULL,
    song TEXT
)