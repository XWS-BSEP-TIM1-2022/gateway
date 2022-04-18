package api

import (
	"context"
	"fmt"
	"gateway/infrastructure/services"
	"gateway/startup/config"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/gateway"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
)

type UserGatewayStruct struct {
	gateway.UnimplementedUserGatewayServer
	config *config.Config
}

func NewUserGateway(c *config.Config) *UserGatewayStruct {
	return &UserGatewayStruct{
		config: c,
	}
}

func (s *UserGatewayStruct) GetRequest(ctx context.Context, in *user.UserIdRequest) (*user.GetResponse, error) {
	userClient := services.NewUserClient(fmt.Sprintf("%s:%s", s.config.UserServiceHost, s.config.UserServicePort))
	return userClient.GetRequest(ctx, in)
}

func (s *UserGatewayStruct) GetAllRequest(ctx context.Context, in *user.EmptyRequest) (*user.GetAllUsers, error) {
	userClient := services.NewUserClient(fmt.Sprintf("%s:%s", s.config.UserServiceHost, s.config.UserServicePort))
	return userClient.GetAllRequest(ctx, in)
}

func (s *UserGatewayStruct) PostRequest(ctx context.Context, in *user.UserRequest) (*user.GetResponse, error) {
	userClient := services.NewUserClient(fmt.Sprintf("%s:%s", s.config.UserServiceHost, s.config.UserServicePort))
	return userClient.PostRequest(ctx, in)
}

func (s *UserGatewayStruct) UpdateRequest(ctx context.Context, in *user.UserRequest) (*user.GetResponse, error) {
	userClient := services.NewUserClient(fmt.Sprintf("%s:%s", s.config.UserServiceHost, s.config.UserServicePort))
	return userClient.UpdateRequest(ctx, in)
}

func (s *UserGatewayStruct) DeleteRequest(ctx context.Context, in *user.UserIdRequest) (*user.EmptyRequest, error) {
	userClient := services.NewUserClient(fmt.Sprintf("%s:%s", s.config.UserServiceHost, s.config.UserServicePort))
	return userClient.DeleteRequest(ctx, in)
}
