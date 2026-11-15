package store

import (
	"context"
	"time"
)

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
	rows, err := s.pool.Query(ctx, `
		SELECT id, task_date, title, progress, difficulty, other_skill, note
		FROM tasks
		WHERE task_date = $1
		ORDER BY id`, date)
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

// ListSkills returns every skill, grouped by kind.
func (s *Store) ListSkills(ctx context.Context) ([]Skill, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, kind FROM skills ORDER BY kind, name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []Skill
	for rows.Next() {
		var sk Skill
		if err := rows.Scan(&sk.ID, &sk.Name, &sk.Kind); err != nil {
			return nil, err
		}
		skills = append(skills, sk)
	}
	return skills, rows.Err()
}
