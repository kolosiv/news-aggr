CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    pub_date TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    link TEXT NOT NULL UNIQUE,
    source_name VARCHAR(255) NOT NULL
);