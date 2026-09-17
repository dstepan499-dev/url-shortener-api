package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	ErrNotFound    = errors.New("url not found")
	ErrAliasExists = errors.New("alias already exists")
)

type URlRecord struct {
	ID          int64     `json:"id"`
	Alias       string    `json:"alias"`
	OriginalURL string    `json:"original_url"`
	Clicks      int       `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}

type Repository struct {
	db *sql.DB
}

func New(connString string) (*Repository, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) InitSchema(schemaSQL string) error {
	_, err := r.db.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("failed to init schema: %w", err)
	}

	return nil
}

func (r *Repository) SaveURL(ctx context.Context, alias, originalURl string) error {
	query := `INSERT INTO urls (alias, original_url) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, alias, originalURl)
	if err != nil {
		return fmt.Errorf("failed to save url: %w", err)
	}
	return nil
}

func (r *Repository) GetAndIncrement(ctx context.Context, alias string) (string, error) {
	query := `
	UPDATE urls
	SET clicks = clicks + 1
	WHERE alias = $1
	RETURNING original_url;
	`

	var originalURL string
	err := r.db.QueryRowContext(ctx, query, alias).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("failed to query and increment: %w", err)
	}

	return originalURL, nil
}

func (r *Repository) GetAnalytics(ctx context.Context, alias string) (*URlRecord, error) {
	query := `
	SELECT id, alias, original_url, clicks, created_at
	FROM urls
	WHERE alias = $1
	`

	var record URlRecord
	err := r.db.QueryRowContext(ctx, query, alias).Scan(
		&record.ID,
		&record.Alias,
		&record.OriginalURL,
		&record.Clicks,
		&record.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get analytics: %w", err)
	}

	return &record, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}
