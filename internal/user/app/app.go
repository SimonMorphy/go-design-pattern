package app

import (
	"github.com/SimonMorphy/go-design-pattern/internal/common/metrics"
	cmd "github.com/SimonMorphy/go-design-pattern/internal/user/app/command"
	qry "github.com/SimonMorphy/go-design-pattern/internal/user/app/query"
	"github.com/SimonMorphy/go-design-pattern/internal/user/infrastructure/storage"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Command Command
	Queries Queries
}

type Command struct {
	Create cmd.CreateUserHandler
	Update cmd.UpdateUserHandler
	Delete cmd.DeleteUserHandler
}

type Queries struct {
	Get  qry.GetUserHandler
	List qry.ListUserHandler
}

func NewApplication() Application {
	repository, cleanUp := storage.NewUserRepository()
	defer func() {
		cleanUp()
	}()
	logger := logrus.NewEntry(logrus.StandardLogger())
	todoMetrics := metrics.NewTodoMetrics()
	return Application{
		Command{
			Create: cmd.NewCreateUsrHandler(repository, logger, todoMetrics),
			Update: cmd.NewUpdateHandler(repository, logger, todoMetrics),
			Delete: cmd.NewDeleteUsrHandler(repository, logger, todoMetrics),
		},
		Queries{
			Get:  qry.NewGetUserHandler(repository, logger, todoMetrics),
			List: qry.NewListUserHandler(repository, logger, todoMetrics),
		},
	}
}
