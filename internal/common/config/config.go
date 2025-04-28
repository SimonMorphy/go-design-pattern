package config

import "github.com/sirupsen/logrus"

func init() {
	NewLogrus(WithLevel(logrus.InfoLevel), WithServiceName("gdp"))
	NewViper()
	InitConn()
}
