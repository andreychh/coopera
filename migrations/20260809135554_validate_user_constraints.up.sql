-- SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
-- SPDX-License-Identifier: MIT
SET LOCAL lock_timeout = '1s';
SET LOCAL statement_timeout = '5s';

ALTER TABLE users VALIDATE CONSTRAINT users_username_length;
ALTER TABLE users VALIDATE CONSTRAINT users_username_shape;
ALTER TABLE users VALIDATE CONSTRAINT users_username_needs_introduction;
