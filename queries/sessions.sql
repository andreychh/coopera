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
