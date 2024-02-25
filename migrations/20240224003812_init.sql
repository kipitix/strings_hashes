-- +goose Up
-- +goose StatementBegin
CREATE TABLE string_hashes
(
    id SERIAL PRIMARY KEY,
    hash TEXT UNIQUE NOT NULL 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE string_hashes;
-- +goose StatementEnd
