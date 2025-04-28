package user

import (
	domain "github.com/SimonMorphy/go-design-pattern/internal/domain/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DataBase string

const (
	MYSQL  DataBase = "mysql"
	MONGO  DataBase = "mongo"
	MEMORY DataBase = "memory"
)

func NewUserRepository() domain.Repository {
	switch DataBase(viper.GetString("database.use")) {
	case MYSQL:
		return NewMysqlUserRepository()
	case MONGO:
		return NewMongoDBUserRepository()
	case MEMORY:
		return NewMemoryUserRepository()
	default:
		logrus.Panic("No Such DataBase")
		return nil
	}
}
