-- name: UpsertTags :batchone
INSERT INTO tags (id, name, color) VALUES ($1, $2, $3) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id;
