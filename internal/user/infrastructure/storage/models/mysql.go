package models

import (
	"fmt"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const/errors"
	"github.com/SimonMorphy/go-design-pattern/internal/infrastructure/creational"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysql        *Mysql
	mysqlFactory = creational.NewSingletonFactory()
)

type Mysql struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func (m Mysql) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.User, m.Password, m.Host, m.Port, m.Database)
}

func MysqlSupplier() (_ interface{}, err error) {
	var db *gorm.DB
	err = viper.UnmarshalKey("mysql", &mysql)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	db, err = gorm.Open(driver.Open(mysql.DSN()))
	dataBasePool := Builder().
		MaxIdleConn(viper.GetInt("database.maxIdleConn")).
		MaxOpenConn(viper.GetInt("database.maxOpenConn")).
		ConnMaxLifeTime(viper.GetDuration("database.connMaxLifeTime")).
		ConnMaxIdleTime(viper.GetDuration("database.connMaxIdleTime")).
		Build()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	d, err := db.DB()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	d.SetMaxIdleConns(dataBasePool.MaxIdleConn)
	d.SetMaxOpenConns(dataBasePool.MaxOpenConn)
	d.SetConnMaxLifetime(dataBasePool.ConnMaxLifeTime)
	d.SetConnMaxIdleTime(dataBasePool.ConnMaxIdleTime)
	return db, nil
}

func InitMysql() (*gorm.DB, func(string), error) {
	mysqlFactory.Register("mysql", MysqlSupplier)
	db, err := GetMysql()
	if err != nil {
		logrus.Error(err)
		return nil, nil, errors.New(errors.ErrnoInternalServerError)
	}
	return db, mysqlFactory.Clear, nil
}
func GetMysql() (*gorm.DB, error) {
	get, err := mysqlFactory.Get("mysql")
	if err != nil {
		logrus.Error("init mysql... with", err)
		return nil, err
	}
	db, ok := get.(*gorm.DB)
	if !ok {
		logrus.Error(err)
		return nil, errors.New(errors.ErrnoCastError)
	}
	return db, nil
}
