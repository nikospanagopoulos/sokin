-- +goose Up
ALTER TABLE tasks
    ADD COLUMN difficulty  INTEGER NOT NULL DEFAULT 5 CHECK (difficulty BETWEEN 1 AND 10),
    ADD COLUMN other_skill TEXT    NOT NULL DEFAULT '',
    ADD COLUMN note        TEXT    NOT NULL DEFAULT '';

CREATE TABLE skills (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    kind TEXT NOT NULL CHECK (kind IN ('language', 'tool', 'concept'))
);

CREATE TABLE task_skills (
    task_id  BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    skill_id BIGINT NOT NULL REFERENCES skills(id),
    PRIMARY KEY (task_id, skill_id)
);

INSERT INTO skills (name, kind) VALUES
    ('Go', 'language'), ('C', 'language'), ('Python', 'language'), ('Rust', 'language'),
    ('JavaScript', 'language'), ('Bash', 'language'), ('SQL', 'language'), ('HTML/CSS', 'language'),
    ('Docker', 'tool'), ('Kubernetes', 'tool'), ('Git', 'tool'), ('Linux', 'tool'),
    ('PostgreSQL', 'tool'), ('Terraform', 'tool'), ('CI/CD', 'tool'),
    ('Networking', 'concept'), ('Security', 'concept'), ('Cloud', 'concept'),
    ('Algorithms', 'concept'), ('Testing', 'concept');

-- +goose Down
DROP TABLE task_skills;
DROP TABLE skills;
ALTER TABLE tasks
    DROP COLUMN note,
    DROP COLUMN other_skill,
    DROP COLUMN difficulty;