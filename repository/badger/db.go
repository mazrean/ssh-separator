package badger

import (
	"fmt"
	"os"
	"runtime"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type DB struct {
	DB *badger.DB
	userCounter prometheus.Counter
}

func NewDB() (*DB, func(), error) {
	return newDB(os.Getenv("BADGER_DIR"), "webshell")
}

func newDB(dir string, prometheusNameSpace string) (*DB, func(), error) {
	db, err := badger.Open(badger.DefaultOptions(dir))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open db file: %w", err)
	}

	userCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: prometheusNameSpace,
		Name:      "user_counter",
	})

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
		return nil, nil, fmt.Errorf("failed in transaction: %w", err)
	}

	return &DB{
		DB: db,
		userCounter: userCounter,
	}, func(){db.Close()}, nil
}
