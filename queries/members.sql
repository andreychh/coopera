-- name: InsertMember :one
INSERT INTO members (team_id, user_id, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetMember :one
SELECT
    id,
    team_id,
    user_id,
    role,
    points,
    left_at,
    created_at
FROM members
WHERE team_id = $1 AND user_id = $2 AND left_at IS NULL;
