-- name: InsertTeam :one
INSERT INTO teams (name)
VALUES ($1)
RETURNING *;
