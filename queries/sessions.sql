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

-- name: EndSession :exec
-- Closes a session. Nothing is deleted: the row stays and carries the
-- moment it ended, because a journal records what happened rather than
-- tidies it away.
--
-- The keys to this session are not touched, and need not be. A refresh
-- token buys keys to a session, and every question about one is asked
-- through the session behind it; from here on the answer is that it is
-- closed. Crossing the keys out separately would put the same fact in
-- two places, and the two would part company at the first sign-out that
-- raced a renewal — the key written a moment later would carry no mark
-- and look alive.
--
-- Two conditions guard the write, and neither reports anything. The
-- person is named so that this can only ever close a session of theirs,
-- whoever calls it and however wrongly. And a session already closed is
-- left alone, so that the moment recorded stays the moment it ended
-- rather than the moment somebody asked again.
--
-- How many rows were touched is not asked, because both answers are the
-- same success: someone who wants out is out either way.
UPDATE sessions
SET revoked_at = NOW()
WHERE
    id = SQLC.ARG(session_id)  -- noqa: RF02
    AND user_id = SQLC.ARG(user_id)  -- noqa: RF02
    AND revoked_at IS NULL;

-- name: InsertSession :one
-- Opens a session. It holds no keys of its own: those live beside it and
-- are replaced many times over, while the session lasts from signing in
-- until signing out. It is the session, not any key, that "sign out"
-- ends and that a list of devices would show.
INSERT INTO sessions (user_id)
VALUES ($1)
RETURNING id;

-- name: SpendRefreshToken :one
-- Spends a key and says whose session it opened. Nothing comes back
-- unless the key was good in every way at once: known, unspent, not yet
-- expired, and belonging to a session still open.
--
-- The session is checked here rather than anywhere else because renewal
-- is the one door the gate does not stand in front of — an access token
-- is expired at exactly the moment renewal is wanted, so there is
-- nothing for the gate to read. Without this join a session closed by
-- signing out would go on renewing itself forever.
--
-- The conditions sit in the statement that writes, not in a look taken
-- beforehand, and that is what stops two callers from spending one key.
-- The second blocks on the row; waking, Postgres re-reads it and checks
-- these conditions again, finds consumed_at set, and updates nothing.
--
-- Why nothing came back is not worked out here. Only the failing path
-- needs to know, and only one of the reasons calls for anything to be
-- done about it.
UPDATE refresh_tokens AS t
SET consumed_at = NOW()
FROM sessions AS s
WHERE
    t.hash = SQLC.ARG(hash)  -- noqa: RF02
    AND t.session_id = s.id
    AND t.consumed_at IS NULL
    AND t.expires_at > NOW()
    AND s.revoked_at IS NULL
RETURNING t.session_id, s.user_id;

-- name: GetSpentRefreshToken :one
-- Asks whether a key that would not open anything is one already spent,
-- and if so, whose session it belonged to.
--
-- A key is spent once and replaced; a second showing means a copy of it
-- exists somewhere. Which of the two hands is the rightful one cannot be
-- known from here, so nothing is guessed — the session ends for both.
--
-- Expiry and closure are deliberately not asked about. Both are already
-- past helping, and neither says anything the holder can act on.
SELECT
    t.session_id,
    s.user_id
FROM refresh_tokens AS t
INNER JOIN sessions AS s ON t.session_id = s.id
WHERE
    t.hash = SQLC.ARG(hash)  -- noqa: RF02
    AND t.consumed_at IS NOT NULL;

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
