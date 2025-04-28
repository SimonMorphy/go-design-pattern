package models

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	once  sync.Once
	db    *gorm.DB
	mysql *Mysql
)

// 饿汉式单例
//func init() {
//	err := viper.UnmarshalKey("mysql", mysql)
//	if err != nil {
//		return
//	}
//	db, err = gorm.Open(driver.Open(mysql.DSN()))
//	if err != nil {
//		return
//	}
//}

type Mysql struct {
	Host     string
	Port     string
	User     string
	Password string
}

func (m Mysql) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/dbname?charset=utf8mb4&parseTime=True&loc=Local", m.User, m.Password, m.Host, m.Port)
}

// GetDB 懒汉单例
func GetDB() (*gorm.DB, error) {
	var err error
	err = viper.UnmarshalKey("mysql", mysql)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if db == nil {
		once.Do(func() {
			db, err = gorm.Open(driver.Open(mysql.DSN()))
			if err != nil {
				logrus.Panic(err)
			}
			dataBasePool := Builder().
				MaxIdleConn(viper.GetInt("database.maxIdleConn")).
				MaxOpenConn(viper.GetInt("database.maxOpenConn")).
				ConnMaxLifeTime(viper.GetDuration("database.connMaxLifeTime")).
				ConnMaxIdleTime(viper.GetDuration("database.connMaxIdleTime")).
				Build()
			d, err := db.DB()
			d.SetMaxIdleConns(dataBasePool.MaxIdleConn)
			d.SetMaxOpenConns(dataBasePool.MaxOpenConn)
			d.SetConnMaxLifetime(dataBasePool.ConnMaxLifeTime)
			d.SetConnMaxIdleTime(dataBasePool.ConnMaxIdleTime)
			if err != nil {
				logrus.Panic(err)
			}
		})
	}
	return db, nil
}
