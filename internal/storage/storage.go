package storage

import "errors"

type AliasUrl struct {
	Alias string `json:"alias"`
	URL   string `json:"url"`
}

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
	ErrAliasExists = errors.New("alias exists")
)
