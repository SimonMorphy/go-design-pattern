package user

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key uint) (*Usr, error)
	Set(ctx context.Context, key uint, value *Usr, expire time.Duration) error
	Delete(ctx context.Context, key uint) error
}
