-- name: CreateProxy :exec
INSERT INTO proxies (id, name, description, schema, host, port, username, password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateProxyTag :copyfrom
INSERT INTO proxy_tag (proxy_id, tag_id) VALUES ($1, $2);

-- name: GetProxyWithTags :many
SELECT sqlc.embed(proxies), tags.id, tags.name, tags.color FROM proxies
LEFT JOIN proxy_tag ON proxies.id = proxy_tag.proxy_id
LEFT JOIN tags ON tags.id = proxy_tag.tag_id
WHERE proxies.id = $1;

-- name: UpdateProxy :exec
UPDATE proxies SET
    name=coalesce(sqlc.narg(name), name),
    description=coalesce(sqlc.narg(description), description),
    schema=coalesce(sqlc.narg(schema), schema),
    host=coalesce(sqlc.narg(host), host),
    port=coalesce(sqlc.narg(port), port),
    username=coalesce(sqlc.narg(username), username),
    password=coalesce(sqlc.narg(password), password)
WHERE id = $1;

-- name: SelectProxyForNoKeyUpdate :one
SELECT id FROM proxies WHERE id = $1 FOR NO KEY UPDATE;

-- name: DeleteProxyTags :exec
DELETE FROM proxy_tag WHERE proxy_id = $1;

-- name: DeleteProxy :one
DELETE FROM proxies WHERE id = $1 RETURNING id;
