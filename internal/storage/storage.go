package storage

import "errors"

var ErrURLNotFound = errors.New("URL not found")

type URLRepository interface {
	SaveURLShort(longURL string) (string, error)
	GetLongURL(code string) (string, error)
}
