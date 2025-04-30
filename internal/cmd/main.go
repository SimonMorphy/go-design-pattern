package main

import (
	"github.com/SimonMorphy/go-design-pattern/internal/common"
	_ "github.com/SimonMorphy/go-design-pattern/internal/common/config"
	"github.com/SimonMorphy/go-design-pattern/internal/common/middleware"
	users "github.com/SimonMorphy/go-design-pattern/internal/ports"
	"github.com/SimonMorphy/go-design-pattern/internal/user/app"
	"github.com/labstack/echo/v4"

	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/swag"
)

func main() {
	e := echo.New()

	// 配置日志中间件
	e.Use(middleware.LogrusMiddleware())

	// 配置 recover 中间件
	e.Use(middleware.RecoverMiddleware())

	// 注册Swagger UI路由
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(
		echoSwagger.URL("/swagger/user.yml"),
		echoSwagger.DocExpansion("list"),
		echoSwagger.PersistAuthorization(true)))

	// 提供OpenAPI规范文件的静态文件服务
	e.Static("/swagger", "./api/openapi")

	// 确保user.yml文件可以被访问
	e.GET("/swagger/user.yml", func(c echo.Context) error {
		return c.File("./api/openapi/user.yml")
	})

	application := app.NewApplication()
	server := NewHttpServer(application, common.BaseResponse{})
	users.RegisterHandlers(e, server)
	err := e.Start(":8080")
	if err != nil {
		logrus.Panic(err)
	}
}
