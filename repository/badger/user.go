package badger

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/dgraph-io/badger/v3"
	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
	"github.com/mazrean/separated-webshell/repository"
)

type User struct {
	db *DB
}

func NewUser(db *DB) *User {
	return &User{
		db: db,
	}
}

func (u *User) Create(ctx context.Context, user *domain.User) error {
	txn, err := getTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}
	if txn == nil {
		return errors.New("no transaction")
	}

	_, err = txn.Get([]byte(user.GetName()))
	if err == nil || !errors.Is(err, badger.ErrKeyNotFound) {
		return repository.ErrUserExist
	}

	if user.HashedPassword == "" {
		return repository.ErrUserPasswordEmpty
	}

	err = txn.Set([]byte(user.GetName()), []byte(user.HashedPassword))
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}
	u.db.userCounter.Inc()

	return nil
}

func (*User) GetPassword(ctx context.Context, userName values.UserName) (values.HashedPassword, error) {
	txn, err := getTransaction(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction: %w", err)
	}
	if txn == nil {
		return "", errors.New("no transaction")
	}

	item, err := txn.Get([]byte(userName))
	if err != nil {
		return "", fmt.Errorf("failed to get password: %w", err)
	}

	var password values.HashedPassword
	err = item.Value(func(val []byte) error {
		password, err = values.NewHashedPassword(string(val))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse value: %w", err)
	}

	return password, nil
}

func (*User) GetAllUser(ctx context.Context) ([]values.UserName, error) {
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

	users := make([]values.UserName, 0, 10)

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		k := item.Key()

		userName, err := values.NewUserName(string(k))
		if err != nil {
			return nil, fmt.Errorf("failed in UserName constructor: %w", err)
		}

		users = append(users, userName)
	}

	return users, nil
}
