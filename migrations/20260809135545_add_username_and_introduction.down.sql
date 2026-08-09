-- SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
-- SPDX-License-Identifier: MIT
SET LOCAL lock_timeout = '1s';
SET LOCAL statement_timeout = '5s';

ALTER TABLE users
DROP CONSTRAINT users_username_length,
DROP CONSTRAINT users_username_shape,
DROP CONSTRAINT users_username_needs_introduction;

ALTER TABLE users
DROP COLUMN username,
DROP COLUMN introduced_at;
