package user

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *Usr) (uint, error)
	List(ctx context.Context, off, lim int) ([]*Usr, error)
	Get(ctx context.Context, ID uint) (*Usr, error)
	Delete(ctx context.Context, ID uint) error
	Update(ctx context.Context, usr *Usr, fun func(context.Context, *Usr) (*Usr, error)) error
}

type RepositoryEmptyError struct{}

func (r RepositoryEmptyError) Error() string {
	return "Empty Repository"
}
