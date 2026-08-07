-- name: InsertTeamWithOwner :one
-- A team and its owner come into being together, in one statement, so a
-- team without an owner is not a state anyone can reach: it does not
-- depend on the caller remembering to open a transaction.
--
-- The owner CTE is not referenced by the final SELECT, and does not need
-- to be. Postgres runs every data-modifying CTE exactly once whether or
-- not anything reads from it.
WITH created AS (
    INSERT INTO teams (name)
    VALUES ($1)
    RETURNING id, name, created_at
),

owner AS (
    INSERT INTO members (team_id, user_id, role)
    SELECT
        created.id,
        $2 AS user_id,
        'owner' AS role
    FROM created
)

SELECT
    id,
    name,
    created_at
FROM created;

-- name: GetTeam :one
SELECT
    id,
    name,
    created_at
FROM teams
WHERE id = $1;

-- name: ListTeamsForMember :many
-- The same membership test as GetTeamForMember, on purpose: a team in
-- this list must be one the person can open, and a team missing from it
-- must be one they cannot. Change the condition here and change it
-- there.
--
-- Ordered by when the person joined rather than when the team was made:
-- this is a list about their own history. The team id breaks ties, so
-- two memberships made in the same instant still come back in a stable
-- order.
SELECT
    t.id,
    t.name,
    t.created_at
FROM teams AS t
INNER JOIN members AS m ON t.id = m.team_id
WHERE m.user_id = $1 AND m.left_at IS NULL
ORDER BY m.created_at DESC, t.id ASC;

-- name: GetTeamForMember :one
SELECT
    t.id,
    t.name,
    t.created_at
FROM teams AS t
INNER JOIN members AS m ON t.id = m.team_id
WHERE t.id = $1 AND m.user_id = $2 AND m.left_at IS NULL;
