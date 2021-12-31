-- +migrate Up
CREATE TABLE enemies (
    id              SERIAL PRIMARY KEY,
    enemy_id        TEXT NOT NULL,
    full_name       TEXT NOT NULL,
    email           TEXT NOT NULL,
    rating          REAL NOT NULL,
    last_updated    TIMESTAMP NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS enemies;