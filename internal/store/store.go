package store

import (
	"context"
	"errors"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

type Entry struct {
	ID       int64
	Date     time.Time
	DidToday string
	Learned  string
	Mistakes string
}

// Open connects to PostgreSQL using the POSTGRES_* environment variables.
func Open(ctx context.Context) (*Store, error) {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	name := os.Getenv("POSTGRES_DB")
	if user == "" || pass == "" || name == "" {
		return nil, errors.New("POSTGRES_* variables are not set (run: set -a; source .env; set +a)")
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pass),
		Host:   "127.0.0.1:5432",
		Path:   name,
	}

	pool, err := pgxpool.New(ctx, u.String())
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

// SaveEntry inserts the entry for a day, or updates it if that day already has one.
func (s *Store) SaveEntry(ctx context.Context, date time.Time, did, learned, mistakes string) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO entries (entry_date, did_today, learned, mistakes)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (entry_date) DO UPDATE SET
			did_today  = EXCLUDED.did_today,
			learned    = EXCLUDED.learned,
			mistakes   = EXCLUDED.mistakes,
			updated_at = now()`,
		date, did, learned, mistakes)
	return err
}

// ListEntries returns the most recent entries, newest first.
func (s *Store) ListEntries(ctx context.Context) ([]Entry, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, entry_date, did_today, learned, mistakes
		FROM entries
		ORDER BY entry_date DESC
		LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.ID, &e.Date, &e.DidToday, &e.Learned, &e.Mistakes); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

// GetEntry returns the entry for a day, or nil if that day has none.
func (s *Store) GetEntry(ctx context.Context, date time.Time) (*Entry, error) {
	var e Entry
	err := s.pool.QueryRow(ctx, `
		SELECT id, entry_date, did_today, learned, mistakes
		FROM entries
		WHERE entry_date = $1`, date).
		Scan(&e.ID, &e.Date, &e.DidToday, &e.Learned, &e.Mistakes)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// CountEntries counts the entries between two days, both included.
func (s *Store) CountEntries(ctx context.Context, from, to time.Time) (int, error) {
	var n int
	err := s.pool.QueryRow(ctx,
		`SELECT count(*) FROM entries WHERE entry_date BETWEEN $1 AND $2`,
		from, to).Scan(&n)
	return n, err
}
