package storage

import (
	"github.com/SimonMorphy/go-design-pattern/internal/user/adapter/cache"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Cache string

const memory, redis Cache = "memory", "redis"

func NewUserCache() (domain.Cache, func()) {
	switch Cache(viper.GetString("cache.use")) {
	case memory:
		return cache.NewMemoryUserCache()
	case redis:
		return cache.NewRedisUserCache()
	default:
		logrus.Panic("No such Cache")
	}
	return nil, nil
}
