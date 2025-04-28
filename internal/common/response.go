package common

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type BaseResponse struct {
}

func (base *BaseResponse) Response(ctx echo.Context, err error, data interface{}) {
	if err != nil {
		base.Error(ctx, err)
	} else {
		base.Success(ctx, data)
	}
}

func (base *BaseResponse) Error(ctx echo.Context, err error) {
	er := ctx.JSON(http.StatusOK, response{
		Errno:   2,
		Message: err.Error(),
		Data:    nil,
	})
	if er != nil {
		logrus.Errorf("failed to return Error to %s", ctx.RealIP())
	}
}

func (base *BaseResponse) Success(ctx echo.Context, data interface{}) {
	err := ctx.JSON(http.StatusOK, response{
		Errno:   0,
		Message: "success",
		Data:    data,
	})
	if err != nil {
		logrus.Errorf("failed to return Success to %s", ctx.RealIP())
	}
}

type response struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id"`
}
