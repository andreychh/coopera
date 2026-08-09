-- name: GetUser :one
-- The introduction is read through the username rather than through
-- introduced_at, though either would do: a CHECK keeps the two together,
-- so one nullable column answers both "has this person introduced
-- themselves" and "what are they called". The timestamp stays behind —
-- nothing outside the database asks when it happened.
SELECT
    id,
    username,
    created_at
FROM users
WHERE id = $1;

-- name: IntroduceUser :one
-- Names a person who has no name yet, and refuses to overwrite one that
-- is already there — so introducing oneself twice touches nothing, no
-- matter what happens in parallel.
--
-- Nothing comes back when no row was touched, and that answers only half
-- the question: either no such person, or a person already introduced.
-- Telling those apart is left to a second query, on the failing path
-- alone, rather than paid for on every call.
--
-- The name is not returned because the caller supplied it. What comes
-- back is what only the database knew.
UPDATE users
SET username = $2, introduced_at = NOW()
WHERE id = $1 AND username IS NULL
RETURNING id, created_at;
