
-- +goose Up
-- +goose StatementBegin
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    CONSTRAINT fk_user 
        FOREIGN KEY(user_id) 
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_feed
        FOREIGN KEY(feed_id) 
        REFERENCES feeds(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feed_follows;
-- +goose StatementEnd
