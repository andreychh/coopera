-- name: InsertInviteLinkAsOwner :one
-- Authorization is part of the write rather than a check before it: the
-- inserted row is selected from members, so a caller who doesn't own the
-- team yields no source row and inserts nothing. Check and write cannot
-- observe different states, so an owner removed mid-request can't leave a
-- link behind.
--
-- No rows means one of "no such team", "not a member", "left the team" or
-- "not the owner", and the caller cannot tell which -- that indistinctness
-- is required: an outsider must not learn whether a team exists.
INSERT INTO invite_links (team_id, code, created_by_member_id, expires_at)
SELECT
    m.team_id,
    $2 AS code,
    m.id,
    $3 AS expires_at
FROM members AS m
WHERE
    m.team_id = $1
    AND m.user_id = $4
    AND m.role = 'owner'
    AND m.left_at IS NULL
RETURNING code, use_count, expires_at, revoked_at, created_at;

-- name: ListInviteLinksByTeam :many
-- Returns the facts, not a verdict: what state each link is in follows
-- from expires_at, revoked_at and the current time, and that derivation
-- lives in the domain (NewLinkState) so there is one of it.
SELECT
    code,
    use_count,
    expires_at,
    revoked_at,
    created_at
FROM invite_links
WHERE team_id = $1
ORDER BY created_at;

-- name: GetInviteLinkByCode :one
SELECT
    id,
    team_id,
    expires_at,
    revoked_at
FROM invite_links
WHERE code = $1
-- Locks the row so a concurrent accept/revoke can't commit between this
-- check and the write it gates: whichever caller acquires the lock
-- first is what the other one sees once its own transaction proceeds.
FOR UPDATE;

-- name: IncrementInviteLinkUseCount :exec
UPDATE invite_links
SET use_count = use_count + 1
WHERE id = $1;

-- name: RevokeInviteLink :exec
UPDATE invite_links
SET revoked_at = NOW()
WHERE id = $1;
