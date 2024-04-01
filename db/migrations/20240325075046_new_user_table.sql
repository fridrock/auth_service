-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
 id SERIAL PRIMARY KEY,
 username VARCHAR(255) NOT NULL,
 email varchar(255) NOT NULL,
 hashed_password varchar(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
