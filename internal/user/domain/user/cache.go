package user

import "context"

type Cache interface {
	Get(ctx context.Context, key string) (*Usr, error)
	Set(ctx context.Context, key string, value *Usr) error
	Delete(ctx context.Context, key string) error
}
