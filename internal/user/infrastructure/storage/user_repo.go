package storage

import (
	"github.com/SimonMorphy/go-design-pattern/internal/user/adapter/cache"
	external2 "github.com/SimonMorphy/go-design-pattern/internal/user/adapter/external"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DataBase string
type Cache string

const (
	MYSQL  DataBase = "mysql"
	MONGO  DataBase = "mongo"
	MEMORY DataBase = "memory"
)

// NewUserRepository <-* Impl Here *-> SimpleFactoryPattern
func NewUserRepository() domain.Repository {
	switch DataBase(viper.GetString("database.use")) {
	case MYSQL:
		return external2.NewMysqlUserRepository()
	case MONGO:
		return external2.NewMongoDBUserRepository()
	case MEMORY:
		return cache.NewMemoryUserRepository()
	default:
		logrus.Panic("No Such DataBase")
		return nil
	}
}

// RepositoryFactory <-* Impl Here *-> abstractFactoryPattern
type RepositoryFactory interface {
	create() domain.Repository
}

type MysqlRepositoryFactory struct {
}

func (m MysqlRepositoryFactory) create() domain.Repository {
	return external2.NewMysqlUserRepository()
}

type MongoRepositoryFactory struct {
}

func (m MongoRepositoryFactory) create() domain.Repository {
	return external2.NewMongoDBUserRepository()
}

type MemoryRepositoryFactory struct {
}

func (m MemoryRepositoryFactory) create() domain.Repository {
	return cache.NewMemoryUserRepository()
}
