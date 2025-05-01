package main

import (
	"context"
	"github.com/SimonMorphy/go-design-pattern/internal/common"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const/errors"
	"github.com/SimonMorphy/go-design-pattern/internal/ports"
	"github.com/SimonMorphy/go-design-pattern/internal/user/app"
	"github.com/SimonMorphy/go-design-pattern/internal/user/app/command"
	dto "github.com/SimonMorphy/go-design-pattern/internal/user/app/dto"
	"github.com/SimonMorphy/go-design-pattern/internal/user/app/query"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type HttpServer struct {
	App app.Application
	common.BaseResponse
}

func (h HttpServer) DeleteUser(ctx echo.Context, id users.UserId) error {
	handle, err := h.App.Command.Delete.Handle(ctx.Request().Context(), command.DeleteUser(id))
	if err != nil {
		logrus.Error(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return h.BaseResponse.Error(ctx, err)
	}
	return h.BaseResponse.Success(ctx, handle.(uint))
}

func NewHttpServer(app app.Application, baseResponse common.BaseResponse) *HttpServer {
	return &HttpServer{App: app, BaseResponse: baseResponse}
}

func (h HttpServer) ListUsers(ctx echo.Context, params users.ListUsersParams) error {
	userList, err := h.App.Queries.List.Handle(ctx.Request().Context(), query.ListUser{
		Offset: *params.Page,
		Limit:  *params.Limit,
	})
	if err != nil {
		return h.BaseResponse.Error(ctx, err)
	}
	return h.BaseResponse.Success(ctx, userList)
}

func (h HttpServer) CreateUser(ctx echo.Context) error {
	var user dto.Usr
	err := ctx.Bind(&user)
	if err != nil {
		h.BaseResponse.Error(ctx, errors.NewWithError(errors.ErrnoBindRequestError, err))
	}
	handle, err := h.App.Command.Create.Handle(ctx.Request().Context(), command.CreateUser{
		Usr: &user,
	})
	if err != nil {
		return h.BaseResponse.Error(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, handle)
}

func (h HttpServer) GetUserById(ctx echo.Context, id users.UserId) error {
	handle, err := h.App.Queries.Get.Handle(ctx.Request().Context(), query.GetUser{ID: uint(id)})
	if err != nil {
		return h.BaseResponse.Error(ctx, err)
	}
	return h.BaseResponse.Success(ctx, handle)
}

func (h HttpServer) UpdateUser(ctx echo.Context, id users.UserId) error {
	var user dto.Usr
	err := ctx.Bind(&user)
	if err != nil {
		return h.BaseResponse.Error(ctx, errors.NewWithError(errors.ErrnoBindRequestError, err))
	}
	err = user.Validate()
	if err != nil {
		return h.BaseResponse.Error(ctx, errors.NewWithError(errors.ErrnoRequestValidateError, err))
	}
	user.ID = uint(id)
	_, err = h.App.Queries.Get.Handle(ctx.Request().Context(), query.GetUser{ID: uint(id)})
	if err != nil {
		return h.BaseResponse.Error(ctx, errors.NewWithError(errors.ErrnoUserNotFoundError, err))
	}
	handle, err := h.App.Command.Update.Handle(ctx.Request().Context(), command.UpdateUser{
		Usr: &user,
		Fn: func(_ context.Context, usr *domain.Usr) (*domain.Usr, error) {
			return usr, nil
		},
	})
	if err != nil {
		return h.BaseResponse.Error(ctx, errors.NewWithError(errors.ErrnoUserModifyFailedError, err))
	}
	return h.BaseResponse.Success(ctx, handle)
}
