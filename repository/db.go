package repository

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
)

var (
	db *badger.DB
)

func Setup() (*badger.DB, error) {
	var err error
	db, err = badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		return nil, fmt.Errorf("failed to open db file: %w", err)
	}

	return db, nil
}
