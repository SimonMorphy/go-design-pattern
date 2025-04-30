package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LogrusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			logrus.WithFields(logrus.Fields{
				"method":  c.Request().Method,
				"uri":     c.Request().RequestURI,
				"status":  c.Response().Status,
				"latency": c.Response().Header().Get(echo.HeaderXRequestID),
			}).Info("Request processed")
			return err
		}
	}
}
