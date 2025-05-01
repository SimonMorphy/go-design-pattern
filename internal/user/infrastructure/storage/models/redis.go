package models

import (
	"context"
	"fmt"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const/errors"
	"github.com/SimonMorphy/go-design-pattern/internal/infrastructure/creational"
	driver "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var (
	redis        Redis
	redisFactory = creational.NewSingletonFactory()
)

type Redis struct {
	IP           string        `mapstructure:"ip"`
	Port         string        `mapstructure:"port"`
	User         string        `mapstructure:"user"`
	Password     string        `mapstructure:"password"`
	PoolSize     int           `mapstructure:"pool_size"`
	MaxConn      int           `mapstructure:"max_conn"`
	ConnTimeout  time.Duration `mapstructure:"conn_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

func (r Redis) Config() *driver.Options {
	return &driver.Options{
		Network:         "tcp",
		Username:        r.User,
		Password:        r.Password,
		Addr:            fmt.Sprintf("%s:%s", r.IP, r.Port),
		PoolSize:        r.PoolSize,
		MaxActiveConns:  r.MaxConn,
		ConnMaxLifetime: r.ConnTimeout * time.Millisecond,
		ReadTimeout:     r.ReadTimeout * time.Millisecond,
		WriteTimeout:    r.WriteTimeout * time.Millisecond,
	}
}

func redisSupplier() (interface{}, error) {
	if err := viper.UnmarshalKey("redis", &redis); err != nil {
		logrus.Error(err)
		return nil, errors.NewWithError(errors.ErrnoUnmarshalError, err)
	}
	client := driver.NewClient(redis.Config())
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, errors.NewWithError(errors.ErrnoFailedConnectionError, err)
	}
	return client, nil
}

func InitRedis() (*driver.Client, func(string), error) {
	redisFactory.Register("redis", redisSupplier)
	client, err := GetRedis()
	if err != nil {
		logrus.Error(err)
		return nil, nil, errors.NewWithError(errors.ErrnoInternalServerError, err)
	}
	return client, redisFactory.Clear, nil
}

func GetRedis() (*driver.Client, error) {
	instance, err := redisFactory.Get("redis")
	if err != nil {
		return nil, err
	}
	client, ok := instance.(*driver.Client)
	if !ok {
		return nil, errors.NewWithError(errors.ErrnoCastError, err)
	}
	return client, nil
}
