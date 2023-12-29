package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite3.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS URL (
		id INTEGER PRIMARY KEY,
		user TEXT NOT NULL,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(user, url, alias string) error {
	const op = "storage.sqlite.SaveURL"

	URLIsExist, err := s.UsersURLIsExist(user, url)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if URLIsExist {
		return fmt.Errorf("%s: %w", op, storage.ErrURLExists)
	}

	stmt, err := s.db.Prepare("INSERT INTO url(user, url, alias) VALUES(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(user, url, alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"

	row := s.db.QueryRow("SELECT url FROM URL WHERE alias = ?", alias)

	var url string
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteUrl"

	stmt, err := s.db.Prepare("DELETE FROM URL WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	deletedRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if deletedRows == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}

	return nil
}

func (s *Storage) AliasIsExist(alias string) (bool, error) {
	const op = "storage.sqlite.AliasIsExist"

	row := s.db.QueryRow("SELECT url FROM URL WHERE alias = ?", alias)

	var url string
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, fmt.Errorf("%s: %w", op, storage.ErrAliasExists)
}

func (s *Storage) UsersURLIsExist(user, URL string) (bool, error) {
	const op = "storage.sqlite.URLIsExist"

	row := s.db.QueryRow("SELECT id FROM URL WHERE user = ? AND url = ?", user, URL)

	var url string
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
}

func (s *Storage) GetUserAliasURL(user string) ([]storage.AliasUrl, error) {
	const op = "storage.sqlite.GetUserAliasUrl"

	rows, err := s.db.Query("SELECT alias, url FROM URL WHERE user = ?", user)
	if err != nil {
		return nil, nil
	}

	var result []storage.AliasUrl
	var alias, url string

	for rows.Next() {
		err = rows.Scan(&alias, &url)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		result = append(result, storage.AliasUrl{Alias: alias, URL: url})
	}

	return result, nil
}
