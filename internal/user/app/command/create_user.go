package command

import (
	"context"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const"
	"github.com/SimonMorphy/go-design-pattern/internal/common/decorator"
	"github.com/SimonMorphy/go-design-pattern/internal/user/app/dto"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/sirupsen/logrus"
)

type CreateUser struct {
	Usr *dto.Usr
}

type CreateUserResult struct {
	ID uint
}

type CreateUserHandler decorator.CommandHandler[CreateUser, *CreateUserResult]

type createUserHandler struct {
	repository domain.Repository
	cache      domain.Cache
}

func (c createUserHandler) Handle(ctx context.Context, query CreateUser) (*CreateUserResult, error) {
	usr := query.Usr.ToDomain()
	id, err := c.repository.Create(ctx, usr)
	if err != nil {
		return nil, err
	}
	err = c.cache.Delete(ctx, usr.ID)
	if err != nil {
		return nil, errors.NewWithError(errors.ErrnoCacheDelError, err)
	}
	return &CreateUserResult{
		ID: id,
	}, err
}

func NewCreateUsrHandler(
	repository domain.Repository,
	cache domain.Cache,
	logger *logrus.Entry,
	record decorator.MetricsRecord) CreateUserHandler {
	if repository == nil {
		logrus.Panic(domain.RepositoryEmptyError{})
	}
	return decorator.ApplyHandlerDecorators[CreateUser, *CreateUserResult](
		&createUserHandler{repository: repository, cache: cache},
		logger,
		record,
	)
}
