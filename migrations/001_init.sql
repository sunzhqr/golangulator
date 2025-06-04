-- +goose Up
CREATE TABLE IF NOT EXISTS history_entries (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    expression TEXT NOT NULL,
    result DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS history_entries;