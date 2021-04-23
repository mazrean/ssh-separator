package repository

import (
	"fmt"
	"os"
	"runtime"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	db *badger.DB
)

var userCounter = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: "webshell",
	Name:      "user_counter",
})

func Setup() (*badger.DB, error) {
	dir := os.Getenv("BADGER_DIR")

	var err error
	db, err = badger.Open(badger.DefaultOptions(dir))
	if err != nil {
		return nil, fmt.Errorf("failed to open db file: %w", err)
	}

	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = runtime.NumCPU()
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			userCounter.Inc()
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return db, nil
}
