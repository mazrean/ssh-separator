package repository

import (
	"fmt"
	"os"

	badger "github.com/dgraph-io/badger/v3"
)

var (
	db *badger.DB
)

func Setup() (*badger.DB, error) {
	dir := os.Getenv("BADGER_DIR")

	var err error
	db, err = badger.Open(badger.DefaultOptions(dir))
	if err != nil {
		return nil, fmt.Errorf("failed to open db file: %w", err)
	}

	return db, nil
}
