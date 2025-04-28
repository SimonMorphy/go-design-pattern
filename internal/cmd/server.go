package main

import (
	"github.com/SimonMorphy/go-design-pattern/internal/app"
	"github.com/SimonMorphy/go-design-pattern/internal/common"
	"github.com/SimonMorphy/go-design-pattern/internal/ports"
	"github.com/labstack/echo/v4"
)

type HttpServer struct {
	App app.Application
	common.BaseResponse
}

func NewHttpServer(app app.Application, baseResponse common.BaseResponse) *HttpServer {
	return &HttpServer{App: app, BaseResponse: baseResponse}
}

func (h HttpServer) ListUsers(ctx echo.Context, params users.ListUsersParams) error {
	userList, err := h.App.Queries.List.Handle(ctx, params.Page, params.Limit)
	if err != nil {
		h.BaseResponse.Error(ctx, err)
		return err
	}
	h.BaseResponse.Success(ctx, userList)
	return nil
}

func (h HttpServer) CreateUser(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h HttpServer) GetUserById(ctx echo.Context, id users.UserId) error {
	//TODO implement me
	panic("implement me")
}

func (h HttpServer) UpdateUser(ctx echo.Context, id users.UserId) error {
	//TODO implement me
	panic("implement me")
}
