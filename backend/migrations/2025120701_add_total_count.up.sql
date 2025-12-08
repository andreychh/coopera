BEGIN;

ALTER TABLE coopera.memberships
ADD COLUMN IF NOT EXISTS total_points INTEGER NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_memberships_total_points ON coopera.memberships (total_points);

COMMIT;

