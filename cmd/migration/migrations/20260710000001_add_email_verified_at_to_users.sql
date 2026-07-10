-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN email_verified_at DATETIME(3) NULL AFTER avatar;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN email_verified_at;
-- +goose StatementEnd
