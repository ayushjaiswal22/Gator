-- +goose Up
CREATE TABLE IF NOT EXISTS "feeds" (
    id UUID PRIMARY KEY, 
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL, 
    name TEXT NOT NULL, 
    url TEXT UNIQUE NOT NULL, 
    user_id UUID NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
-- +goose Down
DROP TABLE IF EXISTS feeds;