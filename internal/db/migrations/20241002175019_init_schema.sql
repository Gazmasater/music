-- +goose Up
-- +goose StatementBegin
CREATE TABLE artists (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE song_details (
    id SERIAL PRIMARY KEY,
    artist_id INTEGER REFERENCES artists(id) ON DELETE CASCADE,
    group_name VARCHAR(255) NOT NULL,
    song_name VARCHAR(255) NOT NULL,
    release_date TIMESTAMP,
    text JSONB,
    song_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (song_name, artist_id)
);

-- Индекс на поле artist_id для оптимизации запросов по исполнителям
CREATE INDEX idx_song_details_artist_id ON song_details (artist_id);

-- Индекс на поле release_date для ускорения сортировки по дате выхода
CREATE INDEX idx_song_details_release_date ON song_details (release_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_song_details_artist_id;
DROP INDEX IF EXISTS idx_song_details_release_date;
DROP TABLE song_details;
DROP TABLE artists;
-- +goose StatementEnd
