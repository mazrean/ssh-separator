package repository

import "context"

type ITransaction interface {
	Transaction(context.Context, func(ctx context.Context) error) error
	RTransaction(context.Context, func(ctx context.Context) error) error
}
