-- name: CreateProxyBlock :exec
INSERT INTO proxy_blocks (id, name, description) VALUES ($1, $2, $3);

-- name: CreateProxyBlockTag :copyfrom
INSERT INTO proxy_block_tag (proxy_block_id, tag_id) VALUES ($1, $2);

-- name: GetProxyBlockWithTags :many
SELECT sqlc.embed(proxy_blocks), tags.id, tags.name, tags.color FROM proxy_blocks
LEFT JOIN proxy_block_tag ON proxy_blocks.id=proxy_block_tag.proxy_block_id
LEFT JOIN tags ON tags.id=proxy_block_tag.tag_id
WHERE proxy_blocks.id = $1;

-- name: UpdateProxyBlock :exec
UPDATE proxy_blocks SET
    name=coalesce(sqlc.narg(name), name),
    description=coalesce(sqlc.narg(description), description)
WHERE id = $1;

-- name: SelectProxyBlockForNoKeyUpdate :one
SELECT id FROM proxy_blocks WHERE id = $1 FOR NO KEY UPDATE;

-- name: DeleteProxyBlockTags :exec
DELETE FROM proxy_block_tag WHERE proxy_block_id = $1;

-- name: DeleteProxyBlock :one
DELETE FROM proxy_blocks WHERE id = $1 RETURNING id;

-- name: GetProxyBlockTagsCount :one
SELECT count(tag_id) FROM proxy_block_tag WHERE proxy_block_id = $1;

-- name: GetProxiesByProxyBlockId :many
SELECT proxies.id, proxies.name, proxies.description, proxies.schema, proxies.host, proxies.port, proxies.username, proxies.password
FROM proxies
JOIN proxy_tag ON proxies.id = proxy_tag.proxy_id
JOIN proxy_block_tag ON proxy_tag.tag_id = proxy_block_tag.tag_id
WHERE proxy_block_tag.proxy_block_id = $1
GROUP BY proxies.id, proxies.name, proxies.description, proxies.host, proxies.port, proxies.username, proxies.password
HAVING count(proxies.id) = $2::integer;


