package user

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (*Usr, error)
	Set(ctx context.Context, key string, value *Usr, expire time.Duration) error
	Delete(ctx context.Context, key string) error
}
