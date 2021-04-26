package badger

import "testing"

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
