-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

CREATE TABLE IF NOT EXISTS platform_user(
    id            uuid      DEFAULT gen_random_uuid() PRIMARY KEY,
    email         TEXT      NOT NULL,
    password      TEXT      NOT NULL,
    name          TEXT      NOT NULL,
    gender        TEXT      NOT NULL,
    date_of_birth TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS token(
    id          uuid      DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id     uuid      REFERENCES platform_user(id) NOT NULL,
    value       TEXT      NOT NULL,
    issued_at   TIMESTAMP NOT NULL
);

-- add test data
INSERT INTO platform_user (email, password, name, gender, date_of_birth) VALUES ('admin', 'admin', 'admin', 'male', '1990-01-01 00:00:00');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE token;
DROP TABLE platform_user;
-- +goose StatementEnd
