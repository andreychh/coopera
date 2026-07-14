-- name: InsertMember :one
INSERT INTO members (team_id, user_id, role)
VALUES ($1, $2, $3)
RETURNING *;
