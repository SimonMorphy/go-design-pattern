package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewViper() {
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath("internal/common/config")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error(err)
	}
}
