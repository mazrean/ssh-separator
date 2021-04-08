//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
package repository

import "context"

type ITransaction interface {
	Transaction(context.Context, func(ctx context.Context) error) error
	RTransaction(context.Context, func(ctx context.Context) error) error
}
