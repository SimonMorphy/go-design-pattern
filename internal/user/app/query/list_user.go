package query

import (
	"context"

	"github.com/SimonMorphy/go-design-pattern/internal/common/decorator"
	"github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/sirupsen/logrus"
)

type ListUser struct {
	Offset, Limit int
}

type ListUserResult struct {
	Users []*user.Usr `json:"users"`
}

type ListUserHandler decorator.CommandHandler[ListUser, *ListUserResult]

type listUserHandler struct {
	repository user.Repository
}

func (l listUserHandler) Handle(ctx context.Context, query ListUser) (*ListUserResult, error) {
	list, err := l.repository.List(ctx, query.Offset, query.Limit)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &ListUserResult{
		Users: list,
	}, err
}

func NewListUserHandler(
	repository user.Repository,
	logger *logrus.Entry,
	record decorator.MetricsRecord) ListUserHandler {
	if repository == nil {
		logrus.Panic(user.RepositoryEmptyError{})
	}
	if repository == nil {
		logrus.Panic(user.RepositoryEmptyError{})
	}
	return decorator.ApplyHandlerDecorators[ListUser, *ListUserResult](
		&listUserHandler{repository: repository},
		logger,
		record,
	)
}
