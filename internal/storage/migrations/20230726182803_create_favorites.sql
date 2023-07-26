-- +goose Up
-- +goose StatementBegin
CREATE TABLE favorites
(
    id         serial primary key,
    user_id    int references users (id) on delete cascade    not null,
    article_id int references articles (id) on delete cascade not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE favorites;
-- +goose StatementEnd
