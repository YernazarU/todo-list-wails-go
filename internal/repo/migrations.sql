-- Миграция для таблицы задач (tasks)
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    due_date TIMESTAMPTZ,
    priority SMALLINT NOT NULL DEFAULT 2
);
-- Индекс для сортировки по дате
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
-- Индекс для сортировки по приоритету
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
