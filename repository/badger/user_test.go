package badger

import (
	"context"
	"errors"
	"fmt"
	"testing"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
	ctxManager "github.com/mazrean/separated-webshell/pkg/context"
	"github.com/mazrean/separated-webshell/repository"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t.Parallel()

	t.Run("Create", testCreate)
}

func testCreate(t *testing.T) {
	t.Parallel()
	t.Helper()

	db, close, err := newTestDB("user_create")
	if err != nil {
		t.Errorf("failed to create test db: %w", err)
	}
	defer close()

	user := NewUser(db)

	testUserNames := make([]values.UserName, 0, 4)
	for i := 0; i < 4; i++ {
		testUserName, err := values.NewUserName(fmt.Sprintf("user_%d", i))
		if err != nil {
			t.Errorf("failed to create test user name: %w", err)
		}

		testUserNames = append(testUserNames, testUserName)
	}

	testHashedPassword, err := values.NewHashedPassword("$2a$10$hgaZ4iV9VYb9xHOLF/Bu4utNbulE5kVu0akP3u7.5xo/dh5q2o.YC")
	if err != nil {
		t.Errorf("failed to create test hashed password: %w", err)
	}

	tests := []struct {
		description string
		user *domain.User
		isErr bool
		err error
		noTxn bool
		duplicateUser bool
	}{
		{
			description: "create user with password",
			user: domain.NewUser(testUserNames[0], testHashedPassword),
			err: nil,
		},
		{
			description: "create user with empty password",
			user: domain.NewUser(testUserNames[1], ""),
			isErr: true,
			err: repository.ErrUserPasswordEmpty,
		},
		{
			description: "no transaction",
			user: domain.NewUser(testUserNames[2], testHashedPassword),
			isErr: true,
			noTxn: true,
		},
		{
			description: "user already exists",
			user: domain.NewUser(testUserNames[3], testHashedPassword),
			isErr: true,
			err: repository.ErrUserExist,
			duplicateUser: true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctx := context.Background()

			var txn *badger.Txn
			if !test.noTxn {
				txn = db.DB.NewTransaction(true)
				defer txn.Discard()

				ctx = context.WithValue(ctx, ctxManager.TransactionKey, txn)

				if test.duplicateUser {
					err = txn.Set([]byte(test.user.GetName()), []byte(test.user.HashedPassword))
					if err != nil {
						t.Errorf("failed to create user: %w", err)
					}
				}
			}

			err := user.Create(ctx, test.user)

			if !test.isErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)

				if test.err != nil && !errors.Is(err, test.err) {
					t.Errorf("expected error %+v, got %+v", test.err, err)
				}
			}

			if err == nil {
				item, err := txn.Get([]byte(test.user.GetName()))
				assert.NoError(t, err)

				var password values.HashedPassword
				err = item.Value(func(val []byte) error {
					password, err = values.NewHashedPassword(string(val))
					if err != nil {
						return err
					}

					return nil
				})
				assert.NoError(t, err)

				assert.Equal(t, test.user.HashedPassword, password)

				txn.Commit()

				err = db.DB.View(func(txn *badger.Txn) error {
					item, err := txn.Get([]byte(test.user.GetName()))
					if err != nil {
						return fmt.Errorf("failed to get user: %w", err)
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
						return fmt.Errorf("failed to get user password: %w", err)
					}

					assert.Equal(t, test.user.HashedPassword, password)
					return nil
				})
				assert.NoError(t, err)
			}
		})
	}
}
