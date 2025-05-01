package models

import (
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const"
	"github.com/SimonMorphy/go-design-pattern/internal/infrastructure/creational"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	memoryFactory = creational.NewSingletonFactory()
)

func MemoryCacheSupplier() (interface{}, error) {
	c := cache.New(viper.GetDuration("memory.default_expiration"), viper.GetDuration("memory.cleanup_interval"))
	return c, nil
}

func InitMemoryCache() (*cache.Cache, func(string), error) {
	memoryFactory.Register("memory", MemoryCacheSupplier)
	memoryCache, err := GetMemoryCache()
	if err != nil {
		logrus.Error(errors.NewWithError(errors.ErrnoInternalServerError, err))
	}
	return memoryCache, memoryFactory.Clear, nil
}

func GetMemoryCache() (*cache.Cache, error) {
	m, err := memoryFactory.Get("memory")
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	c, ok := m.(*cache.Cache)
	if !ok {
		return nil, errors.New(errors.ErrnoCastError)
	}
	return c, nil
}
