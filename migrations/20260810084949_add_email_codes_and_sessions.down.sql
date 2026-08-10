-- SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
-- SPDX-License-Identifier: MIT
SET LOCAL lock_timeout = '1s';
SET LOCAL statement_timeout = '5s';

DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS sign_in_codes;

ALTER TABLE users
DROP CONSTRAINT users_email_lowercase,
DROP CONSTRAINT users_email_length,
DROP CONSTRAINT users_email_local_length,
DROP CONSTRAINT users_email_shape,
DROP COLUMN email;

DROP FUNCTION IF EXISTS has_email_shape(TEXT);
