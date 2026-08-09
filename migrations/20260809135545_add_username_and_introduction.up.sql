-- SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
-- SPDX-License-Identifier: MIT
-- squawk-ignore-file disallowed-unique-constraint
-- The unique index covers a column created by this very statement, so it
-- indexes nothing and holds its exclusive lock for an instant. The rule guards
-- against locking a column that already has rows in it.

SET LOCAL lock_timeout = '1s';
SET LOCAL statement_timeout = '5s';

ALTER TABLE users
ADD COLUMN username TEXT UNIQUE,
ADD COLUMN introduced_at TIMESTAMPTZ;

ALTER TABLE users
ADD CONSTRAINT users_username_length CHECK (CHAR_LENGTH(username) BETWEEN 3 AND 32) NOT VALID,
ADD CONSTRAINT users_username_shape CHECK (username ~ '^[a-z0-9]([a-z0-9]|_[a-z0-9])*$') NOT VALID,
ADD CONSTRAINT users_username_needs_introduction
CHECK ((introduced_at IS NOT NULL) = (username IS NOT NULL)) NOT VALID;
