package middleware

import (
	"github.com/labstack/echo/v4"
	echoMid "github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// RecoverMiddleware 返回一个配置好的 recover 中间件
func RecoverMiddleware() echo.MiddlewareFunc {
	return echoMid.RecoverWithConfig(echoMid.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
				"stack": string(stack),
			}).Error("Recovered from panic")
			return nil
		},
	})
}
