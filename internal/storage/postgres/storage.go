package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"example.com/internal/shortener"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func Connect(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Init() error {
	query := `
	CREATE TABLE IF NOT EXISTS links (
		id SERIAL PRIMARY KEY,
		long_url TEXT NOT NULL,
		short_code TEXT UNIQUE
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *Storage) SaveURLShort(longURL string) (string, error) {
	//Creating transaction
	tx, err := s.db.Begin()
	if err != nil {
		return "", err
	}
	//Rollback in case of the problem
	defer tx.Rollback()

	code, err := s.checkForDuplicates(longURL, tx)
	if err == nil {
		return code, tx.Commit()
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("check duplicates error: %w ", err)
	}

	var id int

	err = tx.QueryRow(`INSERT INTO links (long_url) VALUES ($1) RETURNING id`, longURL).Scan(&id)
	if err != nil {
		return "", err
	}

	code = shortener.EncodeBase62(id)

	_, err = tx.Exec(`UPDATE links SET short_code = $1 WHERE id = $2`, code, id)
	if err != nil {
		return "", err
	}

	return code, tx.Commit()
}

func (s *Storage) GetLongURL(code string) (string, error) {
	var longURL string
	query := `SELECT long_url FROM links WHERE short_code = $1`
	err := s.db.QueryRow(query, code).Scan(&longURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func (s *Storage) checkForDuplicates(longURL string, tx *sql.Tx) (string, error) {
	query := `SELECT short_code FROM links WHERE long_url = $1`
	var code string
	err := tx.QueryRow(query, longURL).Scan(&code)
	if err != nil {
		return "", err
	}
	return code, nil
}
