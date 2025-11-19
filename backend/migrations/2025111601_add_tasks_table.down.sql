BEGIN;

DROP TABLE IF EXISTS coopera.tasks;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_type t
        JOIN pg_namespace n ON n.oid = t.typnamespace
        WHERE t.typname = 'task_status'
          AND n.nspname = 'coopera'
    ) THEN
DROP TYPE coopera.task_status;
END IF;
END$$;

COMMIT;
