-- name: GetPass :one
-- Answers what the door needs about a pass being shown: whether it still
-- stands, and what its holder is called.
--
-- The token was already read and its signature found good, so who is
-- asking is not in doubt. What a signature cannot say is whether the
-- session behind it is still open — signing out and catching a theft
-- both end one, and neither can reach a token already handed out.
--
-- Nothing comes back when the pass is no longer good, and the three ways
-- that happens are not told apart: no such person, no such session, or a
-- session already ended. All three mean the same to the one holding it.
--
-- The name comes back rather than a yes-or-no about it, so that what
-- counts as being introduced stays decided in one place, above.
SELECT u.username
FROM users AS u
INNER JOIN sessions AS s ON u.id = s.user_id
-- The two parameters are named rather than numbered: with both tables
-- keyed by an id, $1 and $2 would arrive in Go as ID and ID_2, and
-- nothing at the call site would say which is which. sqlfluff reads what
-- is inside SQLC.ARG as a column and asks for a table in front of it,
-- which is why the two refusals below are silenced.
WHERE
    u.id = SQLC.ARG(user_id)  -- noqa: RF02
    AND s.id = SQLC.ARG(session_id)  -- noqa: RF02
    AND s.revoked_at IS NULL;

-- name: InsertSession :one
-- Opens a session. It holds no keys of its own: those live beside it and
-- are replaced many times over, while the session lasts from signing in
-- until signing out. It is the session, not any key, that "sign out"
-- ends and that a list of devices would show.
INSERT INTO sessions (user_id)
VALUES ($1)
RETURNING id;

-- name: InsertRefreshToken :one
-- Writes a key for a session. What is stored is a hash — the token
-- itself goes to its holder and is kept nowhere, so a copy of this table
-- opens nothing.
--
-- A month is how long disuse is tolerated, and every key gets the same
-- month. Since refreshing writes a new one, a session lives on as long
-- as somebody keeps using it and dies quietly when nobody does; age
-- alone never ends it.
--
-- What comes back is the span rather than the moment, for the reason it
-- always is: whoever holds the key counts from receiving it, and the
-- subtraction belongs to the clock that set the deadline.
INSERT INTO refresh_tokens (session_id, hash, expires_at)
VALUES ($1, $2, NOW() + INTERVAL '30 days')
RETURNING EXTRACT(EPOCH FROM (expires_at - NOW()))::BIGINT AS expires_in;
