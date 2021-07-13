package badger

import (
	"context"
	"testing"
)

func TestTransaction(t *testing.T) {
	t.Parallel()

	t.Run("Transaction", testTransaction)
	t.Run("RTransaction", testRTransaction)
}

func testTransaction(t *testing.T) {
	t.Parallel()
	t.Helper()

	db, close, err := newTestDB("transaction_transaction")
	if err != nil {
		t.Errorf("failed to create test db: %w", err)
	}
	defer close()

	transaction := NewTransaction(db)

	tests := []struct{
		description string
		ctx context.Context
		fn func(context.Context) error
		err error
	}{}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			err := transaction.Transaction(test.ctx, test.fn)

		})
	}
}

func testRTransaction(t *testing.T) {
	t.Parallel()
	t.Helper()

	db, close, err := newTestDB("transaction_rtransaction")
	if err != nil {
		t.Errorf("failed to create test db: %w", err)
	}
	defer close()
}
