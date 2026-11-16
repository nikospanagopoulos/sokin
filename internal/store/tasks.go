package store

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type Task struct {
	ID         int64
	Date       time.Time
	Title      string
	Progress   int
	Difficulty int
	OtherSkill string
	Note       string
}

type Skill struct {
	ID   int64
	Name string
	Kind string
}

// AddTask saves a task and links it to the chosen skills in one transaction:
// either everything is saved, or nothing is.
func (s *Store) AddTask(ctx context.Context, t Task, skillIDs []int64) (int64, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx) // does nothing after a successful Commit

	var id int64
	err = tx.QueryRow(ctx, `
		INSERT INTO tasks (task_date, title, progress, difficulty, other_skill, note)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		t.Date, t.Title, t.Progress, t.Difficulty, t.OtherSkill, t.Note,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	for _, skillID := range skillIDs {
		_, err := tx.Exec(ctx,
			`INSERT INTO task_skills (task_id, skill_id) VALUES ($1, $2)`,
			id, skillID)
		if err != nil {
			return 0, err
		}
	}
	return id, tx.Commit(ctx)
}

// ListTasks returns the tasks of one day, oldest first.
func (s *Store) ListTasks(ctx context.Context, date time.Time) ([]Task, error) {
	return s.ListTasksBetween(ctx, date, date)
}

// ListTasksBetween returns the tasks between two days, both included.
func (s *Store) ListTasksBetween(ctx context.Context, from, to time.Time) ([]Task, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, task_date, title, progress, difficulty, other_skill, note
		FROM tasks
		WHERE task_date BETWEEN $1 AND $2
		ORDER BY task_date, id`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Progress,
			&t.Difficulty, &t.OtherSkill, &t.Note)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// SetProgress changes the progress (0-100) of one task.
func (s *Store) SetProgress(ctx context.Context, id int64, progress int) error {
	tag, err := s.pool.Exec(ctx,
		`UPDATE tasks SET progress = $2, updated_at = now() WHERE id = $1`,
		id, progress)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// MoveTask moves a task to another day.
func (s *Store) MoveTask(ctx context.Context, id int64, date time.Time) error {
	tag, err := s.pool.Exec(ctx,
		`UPDATE tasks SET task_date = $2, updated_at = now() WHERE id = $1`,
		id, date)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteTask deletes a task. Its skill links are removed by ON DELETE CASCADE.
func (s *Store) DeleteTask(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
