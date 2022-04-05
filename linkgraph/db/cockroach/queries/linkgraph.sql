
-- name: UpsertLink :one
INSERT INTO links (url, retrieved_at) VALUES ($1, $2) 
ON CONFLICT (url) DO UPDATE SET retrieved_at=GREATEST(links.retrieved_at, $2)
RETURNING id, retrieved_at;

-- name: FindLink :many
SELECT url, retrieved_at FROM links WHERE id=$1;

-- name: LinksInPartition :many
SELECT id, url, retrieved_at FROM links WHERE id >= $1 AND id < $2 AND retrieved_at < $3;

-- name: UpsertEdge :one
INSERT INTO edges (src, dst, updated_at) VALUES ($1, $2, NOW())
ON CONFLICT (src,dst) DO UPDATE SET updated_at=NOW()
RETURNING id, updated_at;

-- name: EdgesInPartition :many
SELECT id, src, dst, updated_at FROM edges WHERE src >= $1 AND dst < $2 AND updated_at < $3;

-- name: RemoveStaleEdges :exec
DELETE FROM edges WHERE src=$1 AND updated_at < $2;
