package models

import (
	"context"
	"fmt"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const"
	"github.com/SimonMorphy/go-design-pattern/internal/infrastructure/creational"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var (
	mongodb        *MongoDB
	mongodbFactory = creational.NewSingletonFactory()
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

func MongoDBSupplier() (_ interface{}, err error) {
	var client *mongo.Client
	err = viper.UnmarshalKey("mongo", &mongodb)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongodb.DSN()))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return client, nil
}

func InitMongoDB() (*mongo.Client, func(string), error) {
	mongodbFactory.Register("mongodb", MongoDBSupplier)
	client, err := GetMongoDB()
	if err != nil {
		logrus.Error(err)
		return nil, nil, errors.New(errors.ErrnoInternalServerError)
	}
	return client, mongodbFactory.Clear, nil
}

func GetMongoDB() (*mongo.Client, error) {
	get, err := mongodbFactory.Get("mongodb")
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	client, ok := get.(*mongo.Client)
	if !ok {
		logrus.Error("type assertion failed for mongo client")
		return nil, errors.New(errors.ErrnoCastError)
	}
	return client, nil
}
