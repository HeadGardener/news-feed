-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id            serial primary key,
    username      varchar(255) not null,
    email         varchar(255) not null unique,
    password_hash varchar(255) not null,
    send_flag     int          not null default 1,
    last_online   timestamp    not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
