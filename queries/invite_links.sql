-- name: InsertInviteLink :one
INSERT INTO invite_links (team_id, code, created_by_member_id, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;
