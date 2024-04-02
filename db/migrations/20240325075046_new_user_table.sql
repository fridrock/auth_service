-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
 id SERIAL PRIMARY KEY,
 username VARCHAR(255) NOT NULL UNIQUE,
 email varchar(255) NOT NULL UNIQUE,
 hashed_password varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_statuses(
    id SERIAL PRIMARY KEY,
    val VARCHAR(255) NOT NULL UNIQUE
);
INSERT INTO user_statuses(val) VALUES ('UNCONFIRMED'), ('CONFIRMED');
CREATE TABLE IF NOT EXISTS users_statuses(
    user_id INTEGER REFERENCES users(id),

    status_id INTEGER REFERENCES user_statuses(id),

    CONSTRAINT users_statuses_pk PRIMARY KEY(user_id, status_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_statuses;
DROP TABLE IF EXISTS statuses;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
