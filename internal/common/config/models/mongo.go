package models

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
	"time"
)

var (
	client  *mongo.Client
	only    sync.Once
	mongodb MongoDB
)

type MongoDB struct {
	Host       string
	Port       int
	User       string
	PassWord   string
	DataBase   string
	Collection string
}

func (d MongoDB) DSN() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d", d.User, d.PassWord, d.Host, d.Port)
}

func GetClient() (*mongo.Client, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	err := viper.UnmarshalKey("mongo", &mongodb)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	only.Do(func() {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongodb.DSN()))
		if err != nil {
			logrus.Panic(err)
		}
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			logrus.Panic(err)
		}
	})
	return client, nil
}
