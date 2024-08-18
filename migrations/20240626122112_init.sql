-- +goose Up
-- +goose StatementBegin
CREATE TABLE proxies (
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL,
    schema      TEXT NOT NULL,
    host        TEXT NOT NULL,
    port        INTEGER NOT NULL,
    username    TEXT,
    password    TEXT
);

CREATE TABLE tags (
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    color       TEXT
);

CREATE TABLE proxy_tag (
    proxy_id    UUID NOT NULL REFERENCES proxies(id) ON DELETE CASCADE,
    tag_id      UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (proxy_id, tag_id)
);

CREATE TABLE proxy_blocks (
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE proxy_block_tag (
    proxy_block_id  UUID NOT NULL REFERENCES proxy_blocks(id) ON DELETE CASCADE,
    tag_id          UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (proxy_block_id, tag_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE proxy_block_tag;
DROP TABLE proxy_blocks;
DROP TABLE proxy_tag;
DROP TABLE tags;
DROP TABLE proxies;
-- +goose StatementEnd
