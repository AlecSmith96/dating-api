-- +goose Up
-- +goose StatementBegin
ALTER TABLE platform_user ADD COLUMN location_latitude FLOAT, ADD COLUMN location_longitude FLOAT;

-- add test data
UPDATE platform_user SET location_latitude = 51.4545, location_longitude = 2.5879 WHERE email = 'admin';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE platform_user DROP COLUMN location_latitude, DROP COLUMN location_longitude;
-- +goose StatementEnd
