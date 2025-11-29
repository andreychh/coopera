BEGIN;

-- Безопасное создание enum coopera.task_status
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'task_status'
          AND n.nspname = 'coopera'
    ) THEN
CREATE TYPE coopera.task_status AS ENUM ('open', 'assigned', 'in_review', 'completed', 'archived');
END IF;
END$$;

CREATE TABLE IF NOT EXISTS coopera.tasks
(
    id           SERIAL PRIMARY KEY,
    team_id      INTEGER NOT NULL,
    title        VARCHAR(100) NOT NULL,
    description  VARCHAR,
    points       INTEGER,
    status       coopera.task_status NOT NULL DEFAULT 'open',
    assigned_to  INTEGER,
    created_by   INTEGER NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'utc'),

    CONSTRAINT fk_team
        FOREIGN KEY (team_id)
            REFERENCES coopera.teams (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_assigned_to
        FOREIGN KEY (assigned_to)
            REFERENCES coopera.memberships (id)
            ON DELETE SET NULL,

    CONSTRAINT fk_created_by
        FOREIGN KEY (created_by)
            REFERENCES coopera.memberships (id)
            ON DELETE RESTRICT,

    CONSTRAINT uq_team_title UNIQUE (team_id, title)
);

CREATE INDEX IF NOT EXISTS idx_tasks_team_id ON coopera.tasks (team_id);
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_to ON coopera.tasks (assigned_to);
CREATE INDEX IF NOT EXISTS idx_tasks_created_by ON coopera.tasks (created_by);

COMMIT;
