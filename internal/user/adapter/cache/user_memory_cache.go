package cache

import (
	"context"
	"github.com/SimonMorphy/go-design-pattern/internal/common/const/errors"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/SimonMorphy/go-design-pattern/internal/user/infrastructure/storage/models"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"time"
)

type MemoryUserCache struct {
	cache *cache.Cache
}

func NewMemoryUserCache() (*MemoryUserCache, func()) {
	memoryCache, f, err := models.InitMemoryCache()
	if err != nil {
		logrus.Panic(err)
	}
	return &MemoryUserCache{
			cache: memoryCache,
		}, func() {
			f("memory")
		}
}

func (m MemoryUserCache) Get(_ context.Context, key string) (*domain.Usr, error) {
	if usr, b := m.cache.Get(key); b {
		logrus.WithField("cache", "hit").Infof("user:%v", usr)
		_usr, ok := usr.(*domain.Usr)
		if ok {
			return _usr, nil
		}
		logrus.WithField("cache", "miss").Infof("user:%v", usr)
		return nil, errors.New(errors.ErrnoCastError)
	}
	return nil, errors.New(errors.ErrnoUserNotFoundError)
}

func (m MemoryUserCache) Set(_ context.Context, key string, value *domain.Usr, expire time.Duration) error {
	m.cache.Set(key, value, expire)
	logrus.WithField("cache", "set").Infof("user:%+v", value)
	return nil
}

func (m MemoryUserCache) Delete(_ context.Context, key string) error {
	m.cache.Delete(key)
	logrus.WithField("cache", "delete").Infof("user:%+v", key)
	return nil
}
