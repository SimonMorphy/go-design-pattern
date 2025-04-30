package command

import (
	"context"
	"github.com/SimonMorphy/go-design-pattern/internal/common/decorator"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/sirupsen/logrus"
)

type DeleteUser uint

type DeleteUserHandler decorator.CommandHandler[DeleteUser, interface{}]

type deleteUserHandler struct {
	repo domain.Repository
}

func (d deleteUserHandler) Handle(ctx context.Context, query DeleteUser) (interface{}, error) {
	err := d.repo.Delete(ctx, uint(query))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return uint(query), nil
}

func NewDeleteUsrHandler(
	repository domain.Repository,
	logger *logrus.Entry,
	record decorator.MetricsRecord) DeleteUserHandler {
	if repository == nil {
		logrus.Panic(domain.RepositoryEmptyError{})
	}
	return decorator.ApplyHandlerDecorators[DeleteUser, interface{}](
		&deleteUserHandler{repo: repository},
		logger,
		record,
	)
}
