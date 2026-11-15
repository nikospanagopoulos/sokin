-- +goose Up

CREATE TABLE tasks (
    id         BIGSERIAL PRIMARY KEY,
    task_date  DATE NOT NULL,
    title      TEXT NOT NULL,
    progress   INTEGER NOT NULL DEFAULT 0 CHECK (progress BETWEEN 0 AND 100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX tasks_task_date_idx ON tasks (task_date);
-- +goose Down
DROP TABLE tasks;
