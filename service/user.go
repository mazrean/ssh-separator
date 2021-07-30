//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
	"github.com/mazrean/separated-webshell/repository"
	"github.com/mazrean/separated-webshell/store"
	"github.com/mazrean/separated-webshell/workspace"
)

type IUser interface {
	New(ctx context.Context, name values.UserName, password values.Password) error
	ResetContainer(ctx context.Context, userName values.UserName) error
	Auth(ctx context.Context, name values.UserName, password values.Password) (bool, error)
}

type User struct {
	ww workspace.IWorkspace
	sw store.IWorkspace
	ru repository.IUser
	rt repository.ITransaction
}

func NewUser(ww workspace.IWorkspace, sw store.IWorkspace, ru repository.IUser, rt repository.ITransaction) *User {
	return &User{
		ww: ww,
		sw: sw,
		ru: ru,
		rt: rt,
	}
}

var (
	// ErrUserExist user already exists
	ErrUserExist = errors.New("user exist")
	// ErrWorkspaceExist workspace already exists
	ErrWorkspaceExist = errors.New("workspace exist")
)

func (u *User) New(ctx context.Context, name values.UserName, password values.Password) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	newHashedPassword, err := values.NewHashedPassword(string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed in HashedPassword constructor: %w", err)
	}

	user := domain.NewUser(name, newHashedPassword)

	err = u.rt.Transaction(ctx, func(ctx context.Context) error {
		err := u.ru.Create(ctx, user)
		if err != nil {
			return fmt.Errorf("create user error: %w", err)
		}

		workspace, err := u.ww.Create(ctx, user.GetName())
		if err != nil {
			return fmt.Errorf("create workspace error: %w", err)
		}
		err = u.sw.Set(ctx, user.GetName(), workspace)
		// TODO: 起動済みwokspaceがstoreに登録されない可能性あり
		if err != nil {
			return fmt.Errorf("failed to set workspace: %w", err)
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

func (u *User) ResetContainer(ctx context.Context, userName values.UserName) error {
	workspace, err := u.sw.Get(ctx, userName)
	if err != nil {
		return ErrInvalidUser
	}

	workspace, err = u.ww.Recreate(ctx, workspace)
	if err != nil {
		return fmt.Errorf("failed to recreate workspace: %w", err)
	}

	err = u.sw.Set(ctx, userName, workspace)
	if err != nil {
		return fmt.Errorf("failed to set workspace: %w", err)
	}

	return nil
}

var (
	// ErrInvalidUser invalid user
	ErrInvalidUser = errors.New("invalid user")
	// ErrIncorrectPassword incorrect password
	ErrIncorrectPassword = errors.New("incorrect password")
)

func (u *User) Auth(ctx context.Context, name values.UserName, password values.Password) (bool, error) {
	var hashedPassword values.HashedPassword
	err := u.rt.RTransaction(ctx, func(ctx context.Context) error {
		var err error
		hashedPassword, err = u.ru.GetPassword(ctx, name)
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

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, ErrIncorrectPassword
	}
	if err != nil {
		return false, fmt.Errorf("compare hash error: %w", err)
	}

	return true, nil
}
