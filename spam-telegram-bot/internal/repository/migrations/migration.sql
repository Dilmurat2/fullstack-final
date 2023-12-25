CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
    telegram_id INTEGER NOT NULL UNIQUE,
    username VARCHAR(255) not null,
    first_name VARCHAR(255),
    last_name VARCHAR(255)
);
