-- SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
-- SPDX-License-Identifier: MIT
-- squawk-ignore-file disallowed-unique-constraint
-- The unique index covers a column created by this very statement, so it
-- indexes nothing and holds its exclusive lock for an instant. The rule guards
-- against locking a column that already has rows in it.
--
-- squawk-ignore-file adding-required-field
-- Same reason: there are no rows to fill. An address is what a person is
-- known by, so a nullable column would say it is optional and force every
-- caller to handle an absence that cannot happen.

SET LOCAL lock_timeout = '1s';
SET LOCAL statement_timeout = '5s';

-- The pattern is the one from the HTML standard, which defines what a browser
-- accepts in <input type="email"> and calls itself a wilful violation of RFC
-- 5322: the real grammar allows addresses no mail system on the web has ever
-- seen. See https://html.spec.whatwg.org/#valid-e-mail-address.
--
-- Two departures, both towards being stricter: the domain must carry a dot, or
-- the address is one we cannot deliver to; and the local part is a dot-atom,
-- so a dot may not lead, trail or double, which the standard permits and RFC
-- 5322 does not.
--
-- It lives in a function because the pattern is long, changes as a whole, and
-- stands over every column holding an address. The rules around it stay
-- separate named constraints, so a violation says which one gave way — and
-- each carries its table's name, since a constraint belongs to a table rather
-- than to the value it guards.
CREATE FUNCTION HAS_EMAIL_SHAPE(address TEXT) RETURNS BOOLEAN
LANGUAGE sql IMMUTABLE AS $$
SELECT address ~ (
    '^[a-z0-9!#$%&''*+/=?^_`{|}~-]+'
    || '(?:\.[a-z0-9!#$%&''*+/=?^_`{|}~-]+)*@'
    || '[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?'
    || '(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)+$'
);
$$;

ALTER TABLE users
ADD COLUMN email TEXT NOT NULL UNIQUE
CONSTRAINT users_email_lowercase CHECK (email = LOWER(email))
CONSTRAINT users_email_length CHECK (CHAR_LENGTH(email) <= 254)
CONSTRAINT users_email_local_length CHECK (CHAR_LENGTH(SPLIT_PART(email, '@', 1)) <= 64)
CONSTRAINT users_email_shape CHECK (HAS_EMAIL_SHAPE(email));

CREATE TABLE sign_in_codes (
    id UUID PRIMARY KEY DEFAULT UUIDV7(),
    email TEXT NOT NULL,
    code TEXT NOT NULL,
    attempts_left BIGINT NOT NULL DEFAULT 5,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT sign_in_codes_email_lowercase CHECK (email = LOWER(email)),
    CONSTRAINT sign_in_codes_email_length CHECK (CHAR_LENGTH(email) <= 254),
    CONSTRAINT sign_in_codes_email_local_length
    CHECK (CHAR_LENGTH(SPLIT_PART(email, '@', 1)) <= 64),
    CONSTRAINT sign_in_codes_email_shape CHECK (HAS_EMAIL_SHAPE(email)),
    CONSTRAINT sign_in_codes_code_shape CHECK (code ~ '^[0-9]{6}$'),
    CONSTRAINT sign_in_codes_attempts_left_not_negative CHECK (attempts_left >= 0)
);

CREATE INDEX sign_in_codes_email_created_at_idx ON sign_in_codes (email, created_at);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT UUIDV7(),
    user_id UUID NOT NULL REFERENCES users (id),
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX sessions_user_id_idx ON sessions (user_id);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT UUIDV7(),
    session_id UUID NOT NULL REFERENCES sessions (id),
    hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX refresh_tokens_session_id_idx ON refresh_tokens (session_id);
