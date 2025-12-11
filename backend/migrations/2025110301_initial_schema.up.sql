BEGIN;

CREATE SCHEMA IF NOT EXISTS coopera;

CREATE TYPE coopera.team_role AS ENUM ('manager', 'member');

CREATE TABLE IF NOT EXISTS coopera.users
(
    id          SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username    VARCHAR(50) UNIQUE NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE TABLE IF NOT EXISTS coopera.teams
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    created_by INTEGER     NOT NULL
);

CREATE TABLE IF NOT EXISTS coopera.memberships
(
    id         SERIAL PRIMARY KEY,
    team_id    INTEGER           NOT NULL,
    user_id  INTEGER           NOT NULL,
    role       coopera.team_role NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

ALTER TABLE coopera.teams
    ADD CONSTRAINT fk_created_by
        FOREIGN KEY (created_by)
            REFERENCES coopera.users (id)
            ON DELETE RESTRICT;

ALTER TABLE coopera.memberships
    ADD CONSTRAINT fk_team
        FOREIGN KEY (team_id)
            REFERENCES coopera.teams (id)
            ON DELETE CASCADE,

    ADD CONSTRAINT fk_member
        FOREIGN KEY (user_id)
            REFERENCES coopera.users (id)
            ON DELETE CASCADE,

    ADD CONSTRAINT unique_membership
        UNIQUE (team_id, user_id);

CREATE INDEX IF NOT EXISTS idx_teams_created_by ON coopera.teams (created_by);

CREATE INDEX IF NOT EXISTS idx_memberships_user_id ON coopera.memberships (user_id);

COMMIT;
