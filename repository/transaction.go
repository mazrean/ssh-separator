package repository

import "context"

type ITransaction interface {
	Transaction(func(ctx context.Context) error) error
}
