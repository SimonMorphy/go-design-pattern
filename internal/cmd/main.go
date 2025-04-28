package main

import (
	"github.com/SimonMorphy/go-design-pattern/internal/common"
	_ "github.com/SimonMorphy/go-design-pattern/internal/common/config"
	users "github.com/SimonMorphy/go-design-pattern/internal/ports"
	"github.com/SimonMorphy/go-design-pattern/internal/user/app"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	logrus.NewEntry(logrus.StandardLogger())
	application := app.NewApplication()
	server := NewHttpServer(application, common.BaseResponse{})
	users.RegisterHandlers(e, server)
	err := e.Start(":8080")
	if err != nil {
		logrus.Panic(err)
	}
}
