package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/SerjLeo/storage_bot/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const (
	pagesTable = "page"
)

type Storage struct {
	db *sql.DB
}

func New(basePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", basePath)
	if err != nil {
		return nil, errors.Wrap(err, "connecting to sqlite")
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "trying to ping database")
	}
	storage := &Storage{db: db}
	err = storage.init()
	if err != nil {
		return nil, errors.Wrap(err, "initializing database")
	}
	return storage, nil
}

func (s *Storage) init() error {
	q := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY, username TEXT, url TEXT)`, pagesTable)
	_, err := s.db.Exec(q)
	return err
}

func (s *Storage) Save(ctx context.Context, p *models.Page) error {
	q := fmt.Sprintf(`INSERT INTO %s (username, url) VALUES (?, ?)`, pagesTable)
	_, err := s.db.ExecContext(ctx, q, p.UserName, p.URL)
	if err != nil {
		return errors.Wrap(err, "[sqlite] saving page")
	}
	return nil
}

func (s *Storage) Remove(ctx context.Context, p *models.Page) error {
	q := fmt.Sprintf(`DELETE FROM %s WHERE username=? AND url=?`, pagesTable)
	_, err := s.db.ExecContext(ctx, q, p.UserName, p.URL)
	if err != nil {
		return errors.Wrap(err, "[sqlite] deleting page")
	}
	return nil
}

func (s *Storage) Pick(ctx context.Context, username string) (*models.Page, error) {
	q := fmt.Sprintf(`SELECT url FROM %s WHERE username = ? ORDER BY RANDOM() LIMIT 1`, pagesTable)
	row := s.db.QueryRowContext(ctx, q, username)

	var url string
	err := row.Scan(&url)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "[sqlite] scanning results")
	}
	return &models.Page{
		URL:      url,
		UserName: username,
	}, nil
}

func (s *Storage) IsExist(ctx context.Context, p *models.Page) (bool, error) {
	q := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE username=? AND url=?`, pagesTable)
	row := s.db.QueryRowContext(ctx, q, p.UserName, p.URL)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "[sqlite] finding page")
	}
	return count > 0, nil
}
