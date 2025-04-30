package cache

import (
	"context"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const"
	"github.com/SimonMorphy/go-design-pattern/internal/infrastructure/creational"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/sirupsen/logrus"
	"sync"
)

var memoryCacheFactory = creational.NewSingletonFactory()

type MemoryUserCache struct {
	lock *sync.RWMutex
	data map[uint]*domain.Usr
}

func (m MemoryUserCache) Get(_ context.Context, key uint) (*domain.Usr, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	data, ok := m.data[key]
	if !ok {
		return nil, errors.New(errors.ErrnoUserNotFoundError)
	}
	return data, nil
}

func (m MemoryUserCache) Set(_ context.Context, key uint, value *domain.Usr) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data[key] = value
	return nil
}

func (m MemoryUserCache) Delete(_ context.Context, key uint) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.data, key)
	return nil
}

func MemoryCacheSupplier() (interface{}, error) {
	return &MemoryUserCache{lock: &sync.RWMutex{}, data: make(map[uint]*domain.Usr)}, nil
}

func GetMemoryCache() (*MemoryUserCache, error) {
	repo, err := memoryCacheFactory.Get("memory")
	if err != nil {
		logrus.Panic(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return nil, err
	}
	cache, ok := repo.(*MemoryUserCache)
	if !ok {
		logrus.Panic(errors.New(errors.ErrnoCastError))
	}
	return cache, nil
}

func NewMemoryUserCache() (*MemoryUserCache, func()) {
	memoryCacheFactory.Register("memory", MemoryCacheSupplier)
	cache, err := GetMemoryCache()
	if err != nil {
		logrus.Panic(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return nil, nil
	}
	return cache, func() {
		memoryCacheFactory.Clear("memory")
	}
}
