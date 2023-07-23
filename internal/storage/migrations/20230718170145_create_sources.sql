-- +goose Up
-- +goose StatementBegin
CREATE TABLE sources
(
    id         serial primary key,
    name       varchar(255) not null,
    feed_url   varchar(255) not null unique,
    created_at timestamp    not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sources;
-- +goose StatementEnd
