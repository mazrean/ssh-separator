//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/repository"
	"github.com/mazrean/separated-webshell/workspace"
)

type IUser interface {
	New(ctx context.Context, user *domain.User) error
	SSHAuth(ctx context.Context, user *domain.User) (bool, error)
	SSHHandler(ctx context.Context, userName domain.UserName, connection *domain.Connection) error
}

type User struct {
	workspace.IWorkspace
	repository.IUser
	repository.ITransaction
}

func NewUser(w workspace.IWorkspace, ru repository.IUser, t repository.ITransaction) (*User, error) {
	ctx := context.Background()

	var users []domain.UserName
	err := t.RTransaction(ctx, func(ctx context.Context) error {
		var err error
		users, err = ru.GetAllUser(ctx)
		if err != nil {
			return fmt.Errorf("failed to get all user: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed in transaction: %w", err)
	}

	for _, user := range users {
		err := w.Create(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to create container: %w", err)
		}
	}

	return &User{
		IWorkspace:   w,
		IUser:        ru,
		ITransaction: t,
	}, nil
}

var (
	ErrUserExist      = errors.New("user exist")
	ErrWorkspaceExist = errors.New("workspace exist")
)

func (u *User) New(ctx context.Context, user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	user.HashedPassword = domain.HashedPassword(hashedPassword)

	err = u.ITransaction.Transaction(ctx, func(ctx context.Context) error {
		err := u.IUser.Create(ctx, user)
		if err != nil {
			return fmt.Errorf("create user error: %w", err)
		}

		err = u.IWorkspace.Create(ctx, user.GetName())
		if err != nil {
			return fmt.Errorf("create workspace error: %w", err)
		}

		return nil
	})
	if errors.Is(err, repository.ErrUserExist) {
		return ErrUserExist
	}
	if errors.Is(err, workspace.ErrWorkspaceExist) {
		return ErrWorkspaceExist
	}
	if err != nil {
		return fmt.Errorf("failed in transaction: %w", err)
	}

	return nil
}

var (
	ErrInvalidUser       = errors.New("invalid user")
	ErrIncorrectPassword = errors.New("incorrect password")
)

func (u *User) SSHAuth(ctx context.Context, user *domain.User) (bool, error) {
	var hashedPassword domain.HashedPassword
	err := u.ITransaction.RTransaction(ctx, func(ctx context.Context) error {
		var err error
		hashedPassword, err = u.IUser.GetPassword(ctx, user.GetName())
		if errors.Is(err, repository.ErrUserNotExist) {
			return ErrInvalidUser
		}
		if err != nil {
			return fmt.Errorf("get password error: %w", err)
		}

		return nil
	})
	if err != nil {
		return false, fmt.Errorf("failed in transaction: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, ErrIncorrectPassword
	}
	if err != nil {
		return false, fmt.Errorf("compare hash error: %w", err)
	}

	return true, nil
}

func (u *User) SSHHandler(ctx context.Context, userName domain.UserName, connection *domain.Connection) error {
	err := u.IWorkspace.Connect(ctx, userName, connection)
	if err != nil {
		return fmt.Errorf("connect to workspace error: %w", err)
	}

	return nil
}
