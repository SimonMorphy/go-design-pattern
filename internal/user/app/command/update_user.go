package command

import (
	"context"
	"github.com/SimonMorphy/go-design-pattern/internal/common/decorator"
	"github.com/SimonMorphy/go-design-pattern/internal/user/app/dto"
	"github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/sirupsen/logrus"
)

type UpdateUser struct {
	usr *dto.Usr
	fn  func(context context.Context, usr *user.Usr) (*user.Usr, error)
}

type UpdateUserHandler decorator.CommandHandler[UpdateUser, interface{}]

type updateUserHandler struct {
	repository user.Repository
}

func (u updateUserHandler) Handle(ctx context.Context, query UpdateUser) (interface{}, error) {
	err := u.repository.Update(ctx, query.usr.ToDomain(), func(_ context.Context, usr *user.Usr) (*user.Usr, error) {
		return usr, nil
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return nil, nil
}

func NewUpdateHandler(
	repository user.Repository,
	logger *logrus.Entry,
	record decorator.MetricsRecord) UpdateUserHandler {
	if repository == nil {
		logrus.Panic(user.RepositoryEmptyError{})
	}
	return decorator.ApplyHandlerDecorators[UpdateUser, interface{}](
		&updateUserHandler{
			repository: repository,
		},
		logger,
		record,
	)
}
