package storage

import (
	"context"
	"time"
)

type Cache[K, V any] interface {
	Get(ctx context.Context, key K) (V, error)
	Set(ctx context.Context, key K, value V, duration time.Duration) error
	Delete(ctx context.Context, key K) error
}
