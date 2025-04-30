package common

import (
	"github.com/labstack/echo/v4"
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

func (base *BaseResponse) Error(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusOK, response{
		Errno:   2,
		Message: err.Error(),
		Data:    nil,
	})
}

func (base *BaseResponse) Success(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, response{
		Errno:   0,
		Message: "success",
		Data:    data,
	})
}

type response struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id"`
}
