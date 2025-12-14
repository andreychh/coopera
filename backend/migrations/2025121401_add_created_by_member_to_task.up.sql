BEGIN;

-- 1. Добавляем колонку
ALTER TABLE coopera.tasks
    ADD COLUMN IF NOT EXISTS created_by_member_id INTEGER;

-- 2. Заполняем для уже существующих задач (если есть данные)
UPDATE coopera.tasks t
SET created_by_member_id = m.id
    FROM coopera.memberships m
WHERE m.user_id = t.created_by
  AND m.team_id = t.team_id
  AND t.created_by_member_id IS NULL;

-- 3. Делаем поле обязательным
ALTER TABLE coopera.tasks
    ALTER COLUMN created_by_member_id SET NOT NULL;

-- 4. Добавляем FK
ALTER TABLE coopera.tasks
    ADD CONSTRAINT fk_tasks_created_by_member
        FOREIGN KEY (created_by_member_id)
            REFERENCES coopera.memberships (id)
            ON DELETE RESTRICT;

-- 5. Индекс для быстрых выборок
CREATE INDEX IF NOT EXISTS idx_tasks_created_by_member
    ON coopera.tasks (created_by_member_id);

COMMIT;
