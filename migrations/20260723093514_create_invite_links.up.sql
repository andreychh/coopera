-- SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
-- SPDX-License-Identifier: MIT
SET LOCAL lock_timeout = '1s';
SET LOCAL statement_timeout = '5s';

CREATE TABLE invite_links (
    id UUID PRIMARY KEY DEFAULT UUIDV7(),
    team_id UUID NOT NULL REFERENCES teams (id),
    code TEXT NOT NULL UNIQUE,
    created_by_member_id UUID NOT NULL REFERENCES members (id),
    use_count BIGINT NOT NULL DEFAULT 0,
    expires_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX invite_links_team_id_idx ON invite_links (team_id);
CREATE INDEX invite_links_created_by_member_id_idx ON invite_links (created_by_member_id);
