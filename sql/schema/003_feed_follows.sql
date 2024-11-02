-- +goose Up
CREATE TABLE IF NOT EXISTS feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_feed
        FOREIGN KEY (feed_id)
        REFERENCES feeds(id)
        ON DELETE CASCADE,
    CONSTRAINT unique_user_feed_pair UNIQUE (user_id, feed_id)
);
-- +goose Down
DROP TABLE IF EXISTS feed_follows;
