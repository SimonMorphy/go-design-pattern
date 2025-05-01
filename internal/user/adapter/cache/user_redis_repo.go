package cache

import (
	"context"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/SimonMorphy/go-design-pattern/internal/user/infrastructure/storage/models"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"strconv"
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

func (r RedisUserCache) Get(ctx context.Context, key uint) (*domain.Usr, error) {
	_key := strconv.FormatUint(uint64(key), 10)
	data, err := r.cache.Get(ctx, _key).Result()
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

func (r RedisUserCache) Set(ctx context.Context, key uint, value *domain.Usr, expire time.Duration) error {
	_, err := r.cache.Set(ctx, strconv.Itoa(int(key)), value, expire).Result()
	if err != nil {
		return errors.NewWithError(errors.ErrnoCacheSetError, err)
	}
	logrus.WithField("cache", "set").Infof("user:%v", value)
	return nil
}

func (r RedisUserCache) Delete(ctx context.Context, key uint) error {
	_, err := r.cache.Del(ctx, strconv.Itoa(int(key))).Result()
	if err != nil {
		return errors.NewWithError(errors.ErrnoCacheDelError, err)
	}
	logrus.WithField("cache", "delete").Infof("user:%v", key)
	return nil
}
