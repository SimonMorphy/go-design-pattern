package command

import (
	"context"
	"github.com/go-playground/validator/v10"

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
	cache      user.Cache
}

func (u updateUserHandler) Handle(ctx context.Context, query UpdateUser) (interface{}, error) {
	usr := query.Usr
	err := validator.New().Struct(usr)
	if err != nil {
		return nil, errors.NewWithError(errors.ErrnoParameterInputError, err)
	}
	_usr := usr.ToDomain()
	err = u.repository.Update(ctx, _usr, query.Fn)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	err = u.cache.Delete(ctx, _usr.ID)
	if err != nil {
		logrus.Error(errors.NewWithError(errors.ErrnoCacheDelError, err))
	}
	return _usr.ID, nil
}

func NewUpdateHandler(
	repository user.Repository,
	cache user.Cache,
	logger *logrus.Entry,
	record decorator.MetricsRecord) UpdateUserHandler {
	if repository == nil {
		logrus.Panic(user.RepositoryEmptyError{})
	}
	return decorator.ApplyHandlerDecorators[UpdateUser, interface{}](
		&updateUserHandler{
			repository: repository,
			cache:      cache,
		},
		logger,
		record,
	)
}
