package cache

import (
	"context"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/SimonMorphy/go-design-pattern/internal/user/infrastructure/storage/models"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"strconv"
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

func (m MemoryUserCache) Get(_ context.Context, key uint) (*domain.Usr, error) {
	if usr, b := m.cache.Get(strconv.Itoa(int(key))); b {
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

func (m MemoryUserCache) Set(_ context.Context, key uint, value *domain.Usr, expire time.Duration) error {
	m.cache.Set(strconv.Itoa(int(key)), value, expire)
	logrus.WithField("cache", "set").Infof("user:%+v", value)
	return nil
}

func (m MemoryUserCache) Delete(ctx context.Context, key uint) error {
	m.cache.Delete(strconv.Itoa(int(key)))
	logrus.WithField("cache", "delete").Infof("user:%+v", key)
	return nil
}
