package service

import (
	"context"
	"fmt"

	"github.com/mazrean/separated-webshell/domain/values"
	"github.com/mazrean/separated-webshell/repository"
	"github.com/mazrean/separated-webshell/store"
	"github.com/mazrean/separated-webshell/workspace"
)

type Setup struct {
	ww workspace.IWorkspace
	sw store.IWorkspace
	repository.ITransaction
	repository.IUser
}

func NewSetup(w workspace.IWorkspace, sw store.IWorkspace, t repository.ITransaction, u repository.IUser) *Setup {
	return &Setup{
		ww:           w,
		sw:           sw,
		ITransaction: t,
		IUser:        u,
	}
}

func (s *Setup) Setup() error {
	ctx := context.Background()

	var users []values.UserName
	err := s.ITransaction.RTransaction(ctx, func(ctx context.Context) error {
		var err error
		users, err = s.IUser.GetAllUser(ctx)
		if err != nil {
			return fmt.Errorf("failed to get all user: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed in transaction: %w", err)
	}

	for _, user := range users {
		workspace, err := s.ww.Create(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to create workspace: %w", err)
		}

		err = s.sw.Set(ctx, user, workspace)
		if err != nil {
			return fmt.Errorf("failed to set workspace: %w", err)
		}
	}

	return nil
}
