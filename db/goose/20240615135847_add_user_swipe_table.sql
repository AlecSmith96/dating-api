-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_swipe(
    id                  uuid    DEFAULT    gen_random_uuid() PRIMARY KEY,
    owner_user_id       uuid    REFERENCES platform_user(id) NOT NULL,
    swiped_user_id      uuid    REFERENCES platform_user(id) NOT NULL,
    positive_preference BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS user_match(
    id              uuid DEFAULT    gen_random_uuid() PRIMARY KEY,
    owner_user_id   uuid REFERENCES platform_user(id) NOT NULL,
    matched_user_id uuid REFERENCES platform_user(id) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_swipe;
DROP TABLE user_match;
-- +goose StatementEnd
