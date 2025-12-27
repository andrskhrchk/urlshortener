package main

import (
	"database/sql"

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

func (s *Storage) SaveURL(longURL string) (int, error) {
	var id int
	query := `INSERT INTO links (long_url) VALUES ($1) RETURNING id`

	err := s.db.QueryRow(query, longURL).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) AddShortCode(id int, code string) error {
	query := `UPDATE links SET short_code = $1 WHERE id = $2`
	_, err := s.db.Exec(query, code, id)
	return err
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
