package repository

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/dgraph-io/badger/v3"
	"github.com/mazrean/separated-webshell/domain"
)

type User struct{}

func (*User) Create(ctx context.Context, user *domain.User) error {
	txn, err := getTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}
	if txn == nil {
		return errors.New("no transaction")
	}

	_, err = txn.Get([]byte(user.Name))
	if err == nil || !errors.Is(err, badger.ErrKeyNotFound) {
		return ErrUserExist
	}

	err = txn.Set([]byte(user.Name), []byte(user.Password))
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	return nil
}

func (*User) GetPassword(ctx context.Context, userName string) (string, error) {
	txn, err := getTransaction(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction: %w", err)
	}
	if txn == nil {
		return "", errors.New("no transaction")
	}

	password, err := txn.Get([]byte(userName))
	if err != nil {
		return "", fmt.Errorf("failed to get password: %w", err)
	}

	return string(password.Key()), nil
}

func (*User) GetAllUser(ctx context.Context) ([]string, error) {
	txn, err := getTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	if txn == nil {
		return nil, errors.New("no transaction")
	}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = runtime.NumCPU()
	it := txn.NewIterator(opts)
	defer it.Close()

	users := make([]string, 0, 10)

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		k := item.Key()
		users = append(users, string(k))
	}

	return users, nil
}
