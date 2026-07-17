-- name: InsertTeam :one
INSERT INTO teams (name)
VALUES ($1)
RETURNING *;

-- name: GetTeam :one
SELECT
    id,
    name,
    created_at
FROM teams
WHERE id = $1;
