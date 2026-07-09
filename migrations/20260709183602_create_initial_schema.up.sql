-- SPDX-FileCopyrightText: 2025-2026 Andrey Chernykh
-- SPDX-License-Identifier: MIT
SET LOCAL lock_timeout = '1s';
SET LOCAL statement_timeout = '5s';

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT UUIDV7(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE teams (
    id UUID PRIMARY KEY DEFAULT UUIDV7(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE members (
    id UUID PRIMARY KEY DEFAULT UUIDV7(),
    team_id UUID NOT NULL REFERENCES teams (id),
    user_id UUID NOT NULL REFERENCES users (id),
    role TEXT NOT NULL CHECK (role IN ('owner', 'member')),
    points BIGINT NOT NULL DEFAULT 0,
    left_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (team_id, user_id)
);

CREATE INDEX members_user_id_idx ON members (user_id);

CREATE TABLE duties (
    id UUID PRIMARY KEY DEFAULT UUIDV7(),
    team_id UUID NOT NULL REFERENCES teams (id),
    name TEXT NOT NULL,
    description TEXT,
    points BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX duties_team_id_idx ON duties (team_id);

CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT UUIDV7(),
    duty_id UUID NOT NULL REFERENCES duties (id),
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'assigned', 'done', 'cancelled')),
    assignee_member_id UUID REFERENCES members (id),
    created_by_member_id UUID NOT NULL REFERENCES members (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX jobs_duty_id_idx ON jobs (duty_id);
CREATE INDEX jobs_assignee_member_id_idx ON jobs (assignee_member_id);
CREATE INDEX jobs_created_by_member_id_idx ON jobs (created_by_member_id);
CREATE INDEX jobs_status_idx ON jobs (status);

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT UUIDV7(),
    team_id UUID NOT NULL REFERENCES teams (id),
    member_id UUID NOT NULL REFERENCES members (id),
    type TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX events_team_id_created_at_idx ON events (team_id, created_at);
CREATE INDEX events_member_id_idx ON events (member_id);

CREATE FUNCTION EVENTS_FORBID_MUTATION() RETURNS TRIGGER AS $$
begin
    raise exception 'events is append-only: % is not allowed', tg_op;
end;
$$ LANGUAGE plpgsql;

CREATE TRIGGER events_forbid_mutation
BEFORE UPDATE OR DELETE ON events
FOR EACH ROW EXECUTE FUNCTION EVENTS_FORBID_MUTATION();
