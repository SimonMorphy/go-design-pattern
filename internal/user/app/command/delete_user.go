package command

import (
	"context"
	"fmt"
	consts "github.com/SimonMorphy/go-design-pattern/internal/common/const"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const/errors"
	"github.com/SimonMorphy/go-design-pattern/internal/common/decorator"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/sirupsen/logrus"
)

type DeleteUser uint

type DeleteUserHandler decorator.CommandHandler[DeleteUser, interface{}]

type deleteUserHandler struct {
	repo  domain.Repository
	cache domain.Cache
}

func (d deleteUserHandler) Handle(ctx context.Context, query DeleteUser) (interface{}, error) {
	err := d.repo.Delete(ctx, uint(query))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	err = d.cache.Delete(ctx, fmt.Sprintf("%s%d", consts.UserPrefix, query))
	if err != nil {
		logrus.Error(errors.NewWithError(errors.ErrnoCacheDelError, err))
	}
	return uint(query), nil
}

func NewDeleteUsrHandler(
	repository domain.Repository,
	cache domain.Cache,
	logger *logrus.Entry,
	record decorator.MetricsRecord) DeleteUserHandler {
	if repository == nil {
		logrus.Panic(domain.RepositoryEmptyError{})
	}
	return decorator.ApplyHandlerDecorators[DeleteUser, interface{}](
		&deleteUserHandler{repo: repository, cache: cache},
		logger,
		record,
	)
}
