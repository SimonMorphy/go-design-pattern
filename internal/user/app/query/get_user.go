package query

import (
	"context"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const"
	"github.com/SimonMorphy/go-design-pattern/internal/common/decorator"
	"github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/sirupsen/logrus"
)

type GetUser struct {
	ID uint
}

type GetUserResult struct {
	Usr *user.Usr
}

type GetUserHandler decorator.CommandHandler[GetUser, *GetUserResult]

type getUserHandler struct {
	repository, cache user.Repository
}

func (g getUserHandler) Handle(ctx context.Context, query GetUser) (*GetUserResult, error) {
	result, err := g.cache.Get(ctx, query.ID)
	if err == nil {
		return &GetUserResult{
			Usr: result,
		}, nil
	}
	result, err = g.repository.Get(ctx, query.ID)
	if result == nil {
		err = errors.NewWithError(errors.ErrnoUserNotFoundError, err)
	}
	_, err = g.cache.Create(ctx, result)
	if err != nil {
		logrus.Error(err)
	}
	return &GetUserResult{
		Usr: result,
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
