-- +goose Up
-- +goose StatementBegin
CREATE TABLE books (
                       id SERIAL PRIMARY KEY,
                       uuid UUID NOT NULL,
                       title TEXT NOT NULL,
                       description TEXT,
                       author TEXT NOT NULL,
                       isbn TEXT NOT NULL,
                       price INT NOT NULL,
                       created_at TIMESTAMPTZ DEFAULT NOW(),
                       updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE books;
-- +goose StatementEnd
