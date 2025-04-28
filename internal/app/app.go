package app

import (
	"github.com/SimonMorphy/go-design-pattern/internal/app/command"
	"github.com/SimonMorphy/go-design-pattern/internal/app/query"
)

type Application struct {
	Command Command
	Queries Queries
}

type Command struct {
	Create command.CreateUserHandler
	Update command.UpdateUserHandler
}

type Queries struct {
	Get  query.GetUserHandler
	List query.ListUserHandler
}
