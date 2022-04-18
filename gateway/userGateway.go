package main

import (
	"context"
	"fmt"
	"gateway/infrastructure/services"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
)

func (s *server) GetRequest(ctx context.Context, in *user.UserIdRequest) (*user.GetResponse, error) {
	userClient := services.NewUserClient(fmt.Sprintf("%s:%s", "localhost", "8085"))
	return userClient.GetRequest(ctx, in)
}

func (s *server) GetAllRequest(ctx context.Context, in *user.EmptyRequest) (*user.GetAllUsers, error) {
	userClient := services.NewUserClient(fmt.Sprintf("%s:%s", "localhost", "8085"))
	return userClient.GetAllRequest(ctx, in)
}

func (s *server) PostRequest(ctx context.Context, in *user.UserRequest) (*user.GetResponse, error) {
	userClient := services.NewUserClient(fmt.Sprintf("%s:%s", "localhost", "8085"))
	return userClient.PostRequest(ctx, in)
}

func (s *server) UpdateRequest(ctx context.Context, in *user.UserRequest) (*user.GetResponse, error) {
	userClient := services.NewUserClient(fmt.Sprintf("%s:%s", "localhost", "8085"))
	return userClient.UpdateRequest(ctx, in)
}

func (s *server) DeleteRequest(ctx context.Context, in *user.UserIdRequest) (*user.EmptyRequest, error) {
	userClient := services.NewUserClient(fmt.Sprintf("%s:%s", "localhost", "8085"))
	return userClient.DeleteRequest(ctx, in)
}
