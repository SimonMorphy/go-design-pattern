package query

import (
	"context"
	"github.com/SimonMorphy/go-design-pattern/internal/common/decorator"
	domain "github.com/SimonMorphy/go-design-pattern/internal/domain/user"
	"github.com/sirupsen/logrus"
)

type ListUser struct {
	offset, limit int
}

type ListUserResult struct {
	users []*domain.Usr
}

type ListUserHandler decorator.CommandHandler[ListUser, *ListUserResult]

type listUserHandler struct {
	repository domain.Repository
}

func (l listUserHandler) Handle(ctx context.Context, query ListUser) (*ListUserResult, error) {
	list, err := l.repository.List(ctx, query.offset, query.limit)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &ListUserResult{
		users: list,
	}, err
}
