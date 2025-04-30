package user

import "context"

type Cache interface {
	Get(ctx context.Context, key uint) (*Usr, error)
	Set(ctx context.Context, key uint, value *Usr) error
	Delete(ctx context.Context, key uint) error
}
