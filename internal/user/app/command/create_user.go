package command

import (
	"context"
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
	repo domain.Repository
}

func (c createUserHandler) Handle(ctx context.Context, query CreateUser) (*CreateUserResult, error) {
	id, err := c.repo.Create(ctx, query.Usr.ToDomain())
	if err != nil {
		return nil, err
	}
	return &CreateUserResult{
		ID: id,
	}, err
}

func NewCreateUsrHandler(
	repository domain.Repository,
	logger *logrus.Entry,
	record decorator.MetricsRecord) CreateUserHandler {
	if repository == nil {
		logrus.Panic(domain.RepositoryEmptyError{})
	}
	return decorator.ApplyHandlerDecorators[CreateUser, *CreateUserResult](
		&createUserHandler{repo: repository},
		logger,
		record,
	)
}
