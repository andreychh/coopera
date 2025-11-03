BEGIN;

CREATE TYPE team_role AS ENUM ('manager', 'member');

CREATE TABLE users
(
    id          SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE TABLE teams
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    created_by INTEGER     NOT NULL,

    CONSTRAINT fk_created_by
        FOREIGN KEY (created_by)
            REFERENCES users (id)
            ON DELETE RESTRICT
);

CREATE TABLE memberships
(
    id         SERIAL PRIMARY KEY,
    team_id    INTEGER   NOT NULL,
    member_id  INTEGER   NOT NULL,
    role       team_role NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),

    CONSTRAINT fk_team
        FOREIGN KEY (team_id)
            REFERENCES teams (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_member
        FOREIGN KEY (member_id)
            REFERENCES users (id)
            ON DELETE CASCADE,

    CONSTRAINT unique_membership
        UNIQUE (team_id, member_id)
);

COMMIT;
