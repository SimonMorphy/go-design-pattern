package cache

import (
	"context"
	"github.com/SimonMorphy/go-design-pattern/internal/common/const/errors"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/SimonMorphy/go-design-pattern/internal/user/infrastructure/storage/models"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

type RedisUserCache struct {
	cache *redis.Client
}

func NewRedisUserCache() (*RedisUserCache, func()) {
	client, cleanUp, err := models.InitRedis()
	if err != nil {
		logrus.Panic(err)
	}
	return &RedisUserCache{cache: client}, func() {
		cleanUp("redis")
	}
}

func (r RedisUserCache) Get(ctx context.Context, key string) (*domain.Usr, error) {
	data, err := r.cache.Get(ctx, key).Result()
	if err != nil {
		logrus.WithField("cache", "miss").Info()
		return nil, errors.NewWithError(errors.ErrnoCacheGetError, err)
	}
	var user domain.Usr
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, errors.NewWithError(errors.ErrnoUnmarshalError, err)
	}
	logrus.WithField("cache", "hit").Infof("user:%v", user)
	return &user, nil
}

func (r RedisUserCache) Set(ctx context.Context, key string, value *domain.Usr, expire time.Duration) error {
	_, err := r.cache.Set(ctx, key, value, expire).Result()
	if err != nil {
		return errors.NewWithError(errors.ErrnoCacheSetError, err)
	}
	logrus.WithField("cache", "set").Infof("user:%v", value)
	return nil
}

func (r RedisUserCache) Delete(ctx context.Context, key string) error {
	_, err := r.cache.Del(ctx, key).Result()
	if err != nil {
		return errors.NewWithError(errors.ErrnoCacheDelError, err)
	}
	logrus.WithField("cache", "delete").Infof("user:%v", key)
	return nil
}
