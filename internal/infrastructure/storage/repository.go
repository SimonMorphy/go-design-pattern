package storage

import "context"

type Repository[T any] interface {
	Create(ctx context.Context, data *T) (uint, error)
	List(ctx context.Context, off, lim int) ([]*T, error)
	Get(ctx context.Context, ID uint) (*T, error)
	Delete(ctx context.Context, ID uint) error
	Update(ctx context.Context, usr *T, fun func(context.Context, *T) (*T, error)) error
}
