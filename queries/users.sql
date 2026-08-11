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

-- name: InsertUser :one
-- Gives the address an account, or hands back the one it already has.
-- Both are the same request as far as the person is concerned: they ask
-- to be let in, and whether the system has seen them before is its own
-- affair rather than news to be announced.
--
-- ON CONFLICT does nothing of substance — it writes the address over
-- itself — and is here only so that a row comes back either way. DO
-- NOTHING would answer with silence on every sign-in after the first and
-- force a second query to tell "already there" from "just now", a
-- difference nothing above needs.
--
-- Only the id comes back. The name is not asked for here: whether the
-- person has one changes nothing about being let in, and the pass issued
-- to someone nameless is the same pass.
INSERT INTO users (email)
VALUES ($1)
ON CONFLICT (email) DO UPDATE SET email = excluded.email
RETURNING id;

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
