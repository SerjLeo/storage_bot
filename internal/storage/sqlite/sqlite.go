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
	q := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY, username TEXT, url TEXT, seen INTEGER)`, pagesTable)
	_, err := s.db.Exec(q)
	return err
}

func (s *Storage) Save(ctx context.Context, p *models.Page) error {
	q := fmt.Sprintf(`INSERT INTO %s (username, url, seen) VALUES (?, ?, 0)`, pagesTable)
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
	q := fmt.Sprintf(`SELECT url FROM %s WHERE username = ? AND seen=0 ORDER BY RANDOM() LIMIT 1`, pagesTable)
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

func (s *Storage) List(ctx context.Context, username string) ([]*models.Page, error) {
	q := fmt.Sprintf(`SELECT url, seen FROM %s WHERE username=?`, pagesTable)
	rows, err := s.db.QueryContext(ctx, q, username)
	if err != nil {
		return nil, errors.Wrap(err, "[sqlite] query rows")
	}
	defer func() { _ = rows.Close() }()
	if err == sql.ErrNoRows {
		return []*models.Page{}, nil
	}
	var pages []*models.Page
	for rows.Next() {
		page := models.Page{
			UserName: username,
		}
		var seen int
		if err := rows.Scan(&page.URL, &seen); err != nil {
			return pages, errors.Wrap(err, "[sqlite] scanning row")
		}
		if seen == 1 {
			page.Seen = true
		}
		pages = append(pages, &page)
	}
	if err = rows.Err(); err != nil {
		return pages, errors.Wrap(err, "[sqlite] getting list")
	}
	return pages, nil
}

func (s *Storage) MarkAsSeen(ctx context.Context, p *models.Page) error {
	q := fmt.Sprintf(`UPDATE %s SET seen = 1 WHERE username=? AND url=?`, pagesTable)
	_, err := s.db.ExecContext(ctx, q, p.UserName, p.URL)
	if err != nil {
		return errors.Wrap(err, "[sqlite] update page")
	}
	return nil
}

func (s *Storage) DeleteSeen(ctx context.Context, username string) error {
	q := fmt.Sprintf(`DELETE FROM %s WHERE username=? AND seen=1`, pagesTable)
	_, err := s.db.ExecContext(ctx, q, username)
	if err != nil {
		return errors.Wrap(err, "[sqlite] delete seen")
	}
	return nil
}
