-- name: InsertSignInCode :one
-- Writes a code for the address unless it was asked for too soon: not
-- within a minute of the last one, and not a sixth within any sixty
-- minutes. Both limits are read from the codes already stored, so there
-- is no counter to drift out of step with them.
--
-- Checking and writing are one statement, and the address is locked for
-- its duration. One without the other is not enough: a second statement
-- could slip in between check and write, while two concurrent copies of
-- this one would both read four codes and both write a fifth, since
-- neither sees the other's unfinished work. The lock is taken by address,
-- so it never stands in the way of anyone else. It is released when the
-- statement ends.
--
-- No row comes back when a limit stands in the way. How long to wait is
-- not worked out here: only the failing path needs that answer, and it
-- asks for it separately.
--
-- What comes back is how long the code has left, not the moment it dies.
-- The subtraction happens here so that both spans this address is told —
-- this one and the wait for the next code — are measured against the
-- same clock, the one that set the deadline in the first place.
WITH held AS (
    SELECT PG_ADVISORY_XACT_LOCK(HASHTEXTEXTENDED($1, 0)) AS lock
),

recent AS (
    SELECT
        MAX(c.created_at) AS latest,
        COUNT(c.id) AS taken
    -- held is joined for the lock it takes, not for the column it has;
    -- without the join the two CTEs are independent and Postgres may read
    -- the codes before the lock is held.
    FROM held  -- noqa: ST11
    LEFT JOIN sign_in_codes AS c
        ON
            c.email = $1
            AND c.created_at > NOW() - INTERVAL '1 hour'
)

INSERT INTO sign_in_codes (email, code, expires_at)
SELECT
    $1 AS email,
    $2 AS code,
    NOW() + INTERVAL '10 minutes' AS expires_at
FROM recent
WHERE
    (recent.latest IS NULL OR recent.latest <= NOW() - INTERVAL '1 minute')
    AND recent.taken < 5
RETURNING EXTRACT(EPOCH FROM (expires_at - NOW()))::BIGINT AS expires_in;

-- name: SignInCodeRetryAfter :one
-- Says how long the address has to wait, in seconds. It is asked after
-- InsertSignInCode has answered either way, and takes the latest of
-- three moments: a minute from the most recent code, an hour from the
-- oldest of five, and now. GREATEST ignores a NULL, so whichever limit
-- does not apply drops out on its own; the zero keeps the answer from
-- going negative when both have passed, and from being NULL for an
-- address with no codes at all.
SELECT
    CEIL(EXTRACT(EPOCH FROM GREATEST(
        MAX(created_at) + INTERVAL '1 minute' - NOW(),
        CASE
            WHEN COUNT(*) >= 5 THEN MIN(created_at) + INTERVAL '1 hour' - NOW()
        END,
        INTERVAL '0'
    )))::BIGINT AS retry_after
FROM sign_in_codes
WHERE
    email = $1
    AND created_at > NOW() - INTERVAL '1 hour';
