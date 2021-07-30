//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package repository

import (
	"context"
	"errors"

	"github.com/mazrean/separated-webshell/domain"
	"github.com/mazrean/separated-webshell/domain/values"
)

var (
	// ErrUserPasswordEmpty user password is empty
	ErrUserPasswordEmpty = errors.New("password is empty")
	// ErrUserExist user already exists
	ErrUserExist = errors.New("user exist error")
	// ErrUserNotExist user not exists
	ErrUserNotExist = errors.New("user not exist error")
)

type IUser interface {
	Create(ctx context.Context, user *domain.User) error
	GetPassword(ctx context.Context, userName values.UserName) (hashedPassword values.HashedPassword, err error)
	GetAllUser(ctx context.Context) (users []values.UserName, err error)
}
