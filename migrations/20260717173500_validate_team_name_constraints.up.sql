-- SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
-- SPDX-License-Identifier: MIT
SET LOCAL lock_timeout = '1s';
SET LOCAL statement_timeout = '5s';

ALTER TABLE teams VALIDATE CONSTRAINT teams_name_length;
ALTER TABLE teams VALIDATE CONSTRAINT teams_name_trimmed;
ALTER TABLE teams VALIDATE CONSTRAINT teams_name_no_control_chars;
