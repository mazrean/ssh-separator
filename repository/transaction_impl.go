package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

type Transaction struct{}

func NewTransaction() *Transaction {
	return &Transaction{}
}

const (
	transactionKey = "transaction"
)

func (*Transaction) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	err := db.Update(func(txn *badger.Txn) error {
		ctx := context.WithValue(ctx, transactionKey, txn)

		return fn(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed in transaction: %w", err)
	}

	return nil
}

func (*Transaction) RTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	err := db.View(func(txn *badger.Txn) error {
		ctx := context.WithValue(ctx, transactionKey, txn)

		return fn(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed in transaction: %w", err)
	}

	return nil
}

func getTransaction(ctx context.Context) (*badger.Txn, error) {
	iTxn := ctx.Value(transactionKey)
	if iTxn == nil {
		return nil, nil
	}

	txn, ok := iTxn.(*badger.Txn)
	if !ok {
		return nil, errors.New("broken transaction")
	}

	return txn, nil
}
