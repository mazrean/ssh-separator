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
	t.Run("GetPassword", testGetPassword)
	t.Run("GetAllUser", testGetAllUser)
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
		description   string
		user          *domain.User
		isErr         bool
		err           error
		noTxn         bool
		duplicateUser bool
	}{
		{
			description: "create user with password",
			user:        domain.NewUser(testUserNames[0], testHashedPassword),
			err:         nil,
		},
		{
			description: "create user with empty password",
			user:        domain.NewUser(testUserNames[1], ""),
			isErr:       true,
			err:         repository.ErrUserPasswordEmpty,
		},
		{
			description: "no transaction",
			user:        domain.NewUser(testUserNames[2], testHashedPassword),
			isErr:       true,
			noTxn:       true,
		},
		{
			description:   "user already exists",
			user:          domain.NewUser(testUserNames[3], testHashedPassword),
			isErr:         true,
			err:           repository.ErrUserExist,
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

				err = txn.Commit()
				if err != nil {
					t.Errorf("failed to commit transaction: %w", err)
				}

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

func testGetPassword(t *testing.T) {
	t.Parallel()

	db, close, err := newTestDB("user_get_password")
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
		txnType     transactionType
		userName    values.UserName
		password    values.HashedPassword
		sameTxn     bool
		isErr       bool
		err         error
	}{
		{
			description: "read transaction",
			txnType:     read,
			userName:    testUserNames[0],
			password:    testHashedPassword,
		},
		{
			description: "write transaction",
			txnType:     write,
			userName:    testUserNames[1],
			password:    testHashedPassword,
		},
		{
			description: "no transaction",
			txnType:     none,
			userName:    testUserNames[2],
			password:    testHashedPassword,
			isErr:       true,
		},
		{
			description: "write transaction",
			txnType:     write,
			userName:    testUserNames[3],
			password:    testHashedPassword,
			sameTxn:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctx := context.Background()

			if !test.sameTxn {
				err := db.DB.Update(func(txn *badger.Txn) error {
					err := txn.Set([]byte(test.userName), []byte(test.password))
					if err != nil {
						return fmt.Errorf("failed to set user password: %w", err)
					}

					return nil
				})
				if err != nil {
					t.Errorf("failed to set user password: %w", err)
				}
			}

			var txn *badger.Txn
			switch test.txnType {
			case read:
				txn = db.DB.NewTransaction(false)
				defer txn.Discard()
			case write:
				txn = db.DB.NewTransaction(true)
				defer txn.Discard()
			}

			ctx = context.WithValue(ctx, ctxManager.TransactionKey, txn)

			if test.sameTxn {
				err := txn.Set([]byte(test.userName), []byte(test.password))
				if err != nil {
					t.Errorf("failed to set user password: %w", err)
				}
			}

			password, err := user.GetPassword(ctx, test.userName)

			if !test.isErr {
				assert.NoError(t, err)

				assert.Equal(t, test.password, password)
			} else {
				assert.Error(t, err)

				if test.err != nil && !errors.Is(err, test.err) {
					t.Errorf("expected error %+v, got %+v", test.err, err)
				}
			}
		})
	}
}

func testGetAllUser(t *testing.T) {
	t.Parallel()

	db, close, err := newTestDB("user_get_all_user")
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

	err = db.DB.Update(func(txn *badger.Txn) error {
		for _, userName := range testUserNames {
			err := txn.Set([]byte(userName), []byte(testHashedPassword))
			if err != nil {
				t.Errorf("failed to set user password: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		t.Errorf("failed to set user password: %w", err)
	}

	tests := []struct {
		description string
		txnType     transactionType
		users       []values.UserName
		isErr       bool
		err         error
	}{
		{
			description: "read transaction",
			txnType:     read,
			users:       testUserNames,
		},
		{
			description: "write transaction",
			txnType:     read,
			users:       testUserNames,
		},
		{
			description: "no transaction",
			txnType:     none,
			users:       testUserNames,
			isErr:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctx := context.Background()

			var txn *badger.Txn
			switch test.txnType {
			case read:
				txn = db.DB.NewTransaction(false)
				defer txn.Discard()
			case write:
				txn = db.DB.NewTransaction(true)
				defer txn.Discard()
			}

			ctx = context.WithValue(ctx, ctxManager.TransactionKey, txn)

			users, err := user.GetAllUser(ctx)

			if !test.isErr {
				assert.NoError(t, err)

				assert.Equal(t, test.users, users)
			} else {
				assert.Error(t, err)

				if test.err != nil && !errors.Is(err, test.err) {
					t.Errorf("expected error %+v, got %+v", test.err, err)
				}
			}
		})
	}
}
