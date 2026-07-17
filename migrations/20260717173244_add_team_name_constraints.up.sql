-- SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
-- SPDX-License-Identifier: MIT
SET LOCAL lock_timeout = '1s';
SET LOCAL statement_timeout = '5s';

ALTER TABLE teams
ADD CONSTRAINT teams_name_length CHECK (CHAR_LENGTH(name) BETWEEN 1 AND 100) NOT VALID,
ADD CONSTRAINT teams_name_trimmed CHECK (name = BTRIM(name)) NOT VALID,
ADD CONSTRAINT teams_name_no_control_chars CHECK (name !~ '[[:cntrl:]]') NOT VALID;
