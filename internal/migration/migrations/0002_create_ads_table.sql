-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ads (
    id BIGSERIAL PRIMARY KEY,
    user_login TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    image TEXT NOT NULL,
    price BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    CONSTRAINT fk_users FOREIGN KEY (user_login) REFERENCES users(login)
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS ads;
