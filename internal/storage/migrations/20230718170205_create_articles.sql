-- +goose Up
-- +goose StatementBegin
CREATE TABLE articles
(
    id           serial primary key,
    source_id    int references sources (id) on delete cascade not null,
    title        varchar(255)                                  not null,
    link         varchar(255)                                  not null,
    summary      text                                          not null,
    published_at timestamp                                     not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE articles;
-- +goose StatementEnd
