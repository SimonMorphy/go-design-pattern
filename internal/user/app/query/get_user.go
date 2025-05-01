package query

import (
	"context"
	"fmt"
	"github.com/SimonMorphy/go-design-pattern/internal/common/const"
	"github.com/SimonMorphy/go-design-pattern/internal/common/const/errors"
	"github.com/SimonMorphy/go-design-pattern/internal/common/decorator"
	"github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/sirupsen/logrus"
	"time"
)

type GetUser struct {
	ID uint
}

type GetUserResult struct {
	Usr *user.Usr
}

type GetUserHandler decorator.CommandHandler[GetUser, *GetUserResult]

type getUserHandler struct {
	repository user.Repository
	cache      user.Cache
}

func (g getUserHandler) Handle(ctx context.Context, query GetUser) (*GetUserResult, error) {
	result, err := g.cache.Get(ctx, fmt.Sprintf("%s%d", consts.UserPrefix, query.ID))
	if err == nil {
		return &GetUserResult{
			Usr: result,
		}, nil
	}
	result, err = g.repository.Get(ctx, query.ID)
	if result == nil {
		return nil, errors.NewWithError(errors.ErrnoUserNotFoundError, err)
	}
	err = g.cache.Set(ctx, fmt.Sprintf("%s%d", consts.UserPrefix, result.ID), result, time.Hour)
	if err != nil {
		logrus.Error(errors.NewWithError(errors.ErrnoCacheSetError, err))
	}
	return &GetUserResult{
		Usr: result,
	}, nil
}

func NewGetUserHandler(
	repository user.Repository,
	cache user.Cache,
	logger *logrus.Entry,
	record decorator.MetricsRecord) GetUserHandler {
	if repository == nil {
		logrus.Panic(user.RepositoryEmptyError{})
	}
	return decorator.ApplyHandlerDecorators[GetUser, *GetUserResult](
		&getUserHandler{repository: repository, cache: cache},
		logger,
		record,
	)
}
