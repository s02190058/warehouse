package service

import "context"

type Transactor interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}
