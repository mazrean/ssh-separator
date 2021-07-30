package badger

import (
	"context"
	"errors"
	"fmt"
	"testing"

	badger "github.com/dgraph-io/badger/v3"
	ctxManager "github.com/mazrean/separated-webshell/pkg/context"
	"github.com/stretchr/testify/assert"
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

	testErr := errors.New("test error")

	tests := []struct {
		description string
		ctx         context.Context
		fn          func(context.Context) error
		err         error
	}{
		{
			description: "no error(no process)",
			ctx:         context.Background(),
			fn: func(ctx context.Context) error {
				return nil
			},
			err: nil,
		},
		{
			description: "error(error in func)",
			ctx:         context.Background(),
			fn: func(ctx context.Context) error {
				return testErr
			},
			err: testErr,
		},
		{
			description: "no error(get transaction)",
			ctx:         context.Background(),
			fn: func(ctx context.Context) error {
				iTxn := ctx.Value(ctxManager.TransactionKey)
				if iTxn == nil {
					return errors.New("no transaction")
				}

				_, ok := iTxn.(*badger.Txn)
				if !ok {
					return errors.New("invalid transaction")
				}

				return nil
			},
			err: nil,
		},
		{
			description: "no error(getTransaction)",
			ctx:         context.Background(),
			fn: func(ctx context.Context) error {
				_, err := getTransaction(ctx)
				if err != nil {
					return fmt.Errorf("failed to get transaction: %w", err)
				}

				return nil
			},
			err: nil,
		},
		{
			description: "no error(read)",
			ctx:         context.Background(),
			fn: func(ctx context.Context) error {
				txn, err := getTransaction(ctx)
				if err != nil {
					return fmt.Errorf("failed to get transaction: %w", err)
				}

				_, err = txn.Get([]byte("transaction-key"))
				if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
					return fmt.Errorf("failed to get item: %w", err)
				}

				return nil
			},
			err: nil,
		},
		{
			description: "no error(write)",
			ctx:         context.Background(),
			fn: func(ctx context.Context) error {
				txn, err := getTransaction(ctx)
				if err != nil {
					return fmt.Errorf("failed to get transaction: %w", err)
				}

				err = txn.Set([]byte("transaction-key"), []byte("transaction-value"))
				if err != nil {
					return fmt.Errorf("failed to get item: %w", err)
				}

				return nil
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			err := transaction.Transaction(test.ctx, test.fn)
			if test.err == nil {
				assert.NoError(t, err)
			} else {
				if !errors.Is(err, test.err) {
					t.Errorf("expected %+v, got %+v", test.err, err)
				}
			}
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

	transaction := NewTransaction(db)

	testErr := errors.New("test error")

	tests := []struct {
		description string
		ctx         context.Context
		fn          func(context.Context) error
		err         error
	}{
		{
			description: "no error(no process)",
			ctx:         context.Background(),
			fn: func(ctx context.Context) error {
				return nil
			},
			err: nil,
		},
		{
			description: "error(error in func)",
			ctx:         context.Background(),
			fn: func(ctx context.Context) error {
				return testErr
			},
			err: testErr,
		},
		{
			description: "no error(get transaction)",
			ctx:         context.Background(),
			fn: func(ctx context.Context) error {
				iTxn := ctx.Value(ctxManager.TransactionKey)
				if iTxn == nil {
					return errors.New("no transaction")
				}

				_, ok := iTxn.(*badger.Txn)
				if !ok {
					return errors.New("invalid transaction")
				}

				return nil
			},
			err: nil,
		},
		{
			description: "no error(getTransaction)",
			ctx:         context.Background(),
			fn: func(ctx context.Context) error {
				_, err := getTransaction(ctx)
				if err != nil {
					return fmt.Errorf("failed to get transaction: %w", err)
				}

				return nil
			},
			err: nil,
		},
		{
			description: "no error(read)",
			ctx:         context.Background(),
			fn: func(ctx context.Context) error {
				txn, err := getTransaction(ctx)
				if err != nil {
					return fmt.Errorf("failed to get transaction: %w", err)
				}

				_, err = txn.Get([]byte("transaction-key"))
				if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
					return fmt.Errorf("failed to get item: %w", err)
				}

				return nil
			},
			err: nil,
		},
		{
			description: "error(write)",
			ctx:         context.Background(),
			fn: func(ctx context.Context) error {
				txn, err := getTransaction(ctx)
				if err != nil {
					return fmt.Errorf("failed to get transaction: %w", err)
				}

				err = txn.Set([]byte("transaction-key"), []byte("transaction-value"))
				if err != nil {
					return fmt.Errorf("failed to get item: %w", err)
				}

				return nil
			},
			err: badger.ErrReadOnlyTxn,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			err := transaction.RTransaction(test.ctx, test.fn)
			if test.err == nil {
				assert.NoError(t, err)
			} else {
				if !errors.Is(err, test.err) {
					t.Errorf("expected %+v, got %+v", test.err, err)
				}
			}
		})
	}
}
