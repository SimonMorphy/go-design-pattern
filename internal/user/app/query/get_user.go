package query

import (
	"context"
	"github.com/SimonMorphy/go-design-pattern/internal/common/decorator"
	"github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/sirupsen/logrus"
)

type GetUser struct {
	ID uint
}

type GetUserResult struct {
	usr *user.Usr
}

type GetUserHandler decorator.CommandHandler[GetUser, *GetUserResult]

type getUserHandler struct {
	repository user.Repository
}

func (g getUserHandler) Handle(ctx context.Context, query GetUser) (*GetUserResult, error) {
	result, err := g.repository.Get(ctx, query.ID)
	return &GetUserResult{
		usr: result,
	}, err
}

func NewGetUserHandler(
	repository user.Repository,
	logger *logrus.Entry,
	record decorator.MetricsRecord) GetUserHandler {
	if repository == nil {
		logrus.Panic(user.RepositoryEmptyError{})
	}
	return decorator.ApplyHandlerDecorators[GetUser, *GetUserResult](
		&getUserHandler{repository: repository},
		logger,
		record,
	)
}
