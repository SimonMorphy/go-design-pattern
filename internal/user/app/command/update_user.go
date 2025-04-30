package command

import (
	"context"

	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const"
	"github.com/SimonMorphy/go-design-pattern/internal/common/decorator"
	"github.com/SimonMorphy/go-design-pattern/internal/user/app/dto"
	"github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/sirupsen/logrus"
)

type UpdateUser struct {
	Usr *dto.Usr
	Fn  func(context context.Context, usr *user.Usr) (*user.Usr, error)
}

type UpdateUserHandler decorator.CommandHandler[UpdateUser, interface{}]

type updateUserHandler struct {
	repository user.Repository
}

func (u updateUserHandler) Handle(ctx context.Context, query UpdateUser) (interface{}, error) {
	if query.Usr == nil {
		return nil, errors.New(errors.ErrnoParameterInputError)
	}

	domainUser := query.Usr.ToDomain()
	if domainUser == nil {
		return nil, errors.New(errors.ErrnoParameterInputError)
	}

	err := u.repository.Update(ctx, domainUser, query.Fn)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return domainUser.ID, nil
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
