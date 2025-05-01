package storage

import (
	external "github.com/SimonMorphy/go-design-pattern/internal/user/adapter/external"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DataBase string

const (
	MYSQL DataBase = "mysql"
	MONGO DataBase = "mongo"
	JSON  DataBase = "json"
)

// NewUserRepository <-* Impl Here *-> SimpleFactoryPattern
func NewUserRepository() (domain.Repository, func()) {
	switch DataBase(viper.GetString("database.use")) {
	case MYSQL:
		return external.NewMysqlUserRepository()
	case MONGO:
		return external.NewMongoDBUserRepository()
	case JSON:
		return external.NewJsonUserRepository()
	default:
		logrus.Panic("No Such DataBase")
		return nil, nil
	}
}

// RepositoryFactory <-* Impl Here *-> abstractFactoryPattern
type RepositoryFactory interface {
	create() domain.Repository
}

type MysqlRepositoryFactory struct {
}

func (m MysqlRepositoryFactory) create() (domain.Repository, func()) {
	return external.NewMysqlUserRepository()
}

type MongoRepositoryFactory struct {
}

func (m MongoRepositoryFactory) create() (domain.Repository, func()) {
	return external.NewMongoDBUserRepository()
}

type MemoryRepositoryFactory struct {
}

func (m MemoryRepositoryFactory) create() (domain.Repository, func()) {
	panic("TODO")
}
