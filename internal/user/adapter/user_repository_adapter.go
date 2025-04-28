package adapter

import (
	"context"
	"time"
)

type RepositorySyncAdapter interface {
	Load(ctx context.Context) error
	Flush(ctx context.Context) error
	Persist(ctx context.Context) error
	Status(context.Context) (Status, error)
}

type Status struct {
	LastSyncTime      time.Time
	LastTimeUsed      time.Duration
	MemoryDataCount   uint
	ExternalDataCount uint
}
