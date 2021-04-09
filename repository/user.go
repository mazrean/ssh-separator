//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package repository

import (
	"context"
	"errors"

	"github.com/mazrean/separated-webshell/domain"
)

var (
	ErrUserExist    = errors.New("user exist error")
	ErrUserNotExist = errors.New("user not exist error")
)

type IUser interface {
	Create(ctx context.Context, user *domain.User) error
	GetPassword(ctx context.Context, userName domain.UserName) (hashedPassword domain.HashedPassword, err error)
	GetAllUser(ctx context.Context) (users []domain.UserName, err error)
}
