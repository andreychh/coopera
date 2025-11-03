BEGIN;

CREATE SCHEMA IF NOT EXISTS coopera;

CREATE TYPE coopera.team_role AS ENUM ('manager', 'member');

CREATE TABLE IF NOT EXISTS coopera.users
(
    id          SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE TABLE IF NOT EXISTS coopera.teams
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    created_by INTEGER     NOT NULL,

    CONSTRAINT fk_created_by
        FOREIGN KEY (created_by)
            REFERENCES coopera.users (id)
            ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS coopera.memberships
(
    id         SERIAL PRIMARY KEY,
    team_id    INTEGER           NOT NULL,
    member_id  INTEGER           NOT NULL,
    role       coopera.team_role NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),

    CONSTRAINT fk_team
        FOREIGN KEY (team_id)
            REFERENCES coopera.teams (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_member
        FOREIGN KEY (member_id)
            REFERENCES coopera.users (id)
            ON DELETE CASCADE,

    CONSTRAINT unique_membership
        UNIQUE (team_id, member_id)
);

COMMIT;
