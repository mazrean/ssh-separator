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
}

func testRTransaction(t *testing.T) {
	t.Parallel()
	t.Helper()
}
