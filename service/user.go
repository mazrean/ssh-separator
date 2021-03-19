package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/repository"
	"github.com/mazrean/separated-webshell/workspace"
)

type User struct {
	workspace.IWorkspace
	repository.IUser
	repository.ITransaction
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

	user.Password = string(hashedPassword)

	err = u.ITransaction.Transaction(func(ctx context.Context) error {
		err := u.IUser.Create(ctx, user)
		if err != nil {
			return fmt.Errorf("create user error: %w", err)
		}

		err = u.IWorkspace.Create(ctx, user.Name)
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
	hashedPassword, err := u.IUser.GetPassword(ctx, user.Name)
	if errors.Is(err, repository.ErrUserNotExist) {
		return false, ErrInvalidUser
	}
	if err != nil {
		return false, fmt.Errorf("get password error: %w", err)
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

func (u *User) SSHHandler(ctx context.Context, userName string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	err := u.IWorkspace.Connect(ctx, userName, stdin, stdout, stderr)
	if err != nil {
		return fmt.Errorf("connect to workspace error: %w", err)
	}

	return nil
}
