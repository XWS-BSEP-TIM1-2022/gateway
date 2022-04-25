package api

import (
	"context"
	"errors"
	"fmt"
	"gateway/infrastructure/services"
	"gateway/startup/config"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"google.golang.org/grpc/metadata"
)

type UserGatewayStruct struct {
	userService.UnimplementedUserServiceServer
	config     *config.Config
	userClient userService.UserServiceClient
}

func NewUserGateway(c *config.Config) *UserGatewayStruct {
	return &UserGatewayStruct{
		config:     c,
		userClient: services.NewUserClient(fmt.Sprintf("%s:%s", c.UserServiceHost, c.UserServicePort)),
	}
}

func (s *UserGatewayStruct) GetRequest(ctx context.Context, in *user.UserIdRequest) (*user.GetResponse, error) {
	return s.userClient.GetRequest(ctx, in)
}

func (s *UserGatewayStruct) GetAllRequest(ctx context.Context, in *user.EmptyRequest) (*user.UsersResponse, error) {
	return s.userClient.GetAllRequest(ctx, in)
}

func (s *UserGatewayStruct) PostRequest(ctx context.Context, in *user.UserRequest) (*user.GetResponse, error) {
	return s.userClient.PostRequest(ctx, in)
}

func (s *UserGatewayStruct) PostAdminRequest(ctx context.Context, in *user.UserRequest) (*user.GetResponse, error) {
	return s.userClient.PostAdminRequest(ctx, in)
}

func (s *UserGatewayStruct) UpdateRequest(ctx context.Context, in *user.UserRequest) (*user.GetResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	_, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{JwtToken: jwt[0]})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.userClient.UpdateRequest(ctx, in)
}

func (s *UserGatewayStruct) DeleteRequest(ctx context.Context, in *user.UserIdRequest) (*user.EmptyRequest, error) {
	return s.userClient.DeleteRequest(ctx, in)
}

func (s *UserGatewayStruct) LoginRequest(ctx context.Context, in *user.CredentialsRequest) (*user.LoginResponse, error) {
	return s.userClient.LoginRequest(ctx, in)
}

func (s *UserGatewayStruct) SearchUsersRequest(ctx context.Context, in *user.SearchRequest) (*user.UsersResponse, error) {
	return s.userClient.SearchUsersRequest(ctx, in)
}

func (s *UserGatewayStruct) IsUserAuthenticated(ctx context.Context, in *userService.AuthRequest) (*userService.AuthResponse, error) {
	return s.userClient.IsUserAuthenticated(ctx, in)
}

func (s *UserGatewayStruct) UpdatePasswordRequest(ctx context.Context, in *userService.NewPasswordRequest) (*user.GetResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	_, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{JwtToken: jwt[0]})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.userClient.UpdatePasswordRequest(ctx, in)
}

func (s *UserGatewayStruct) GetAllUsersExperienceRequest(ctx context.Context, in *userService.ExperienceRequest) (*user.ExperienceResponse, error) {
	return s.userClient.GetAllUsersExperienceRequest(ctx, in)
}

func (s *UserGatewayStruct) PostExperienceRequest(ctx context.Context, in *user.NewExperienceRequest) (*user.NewExperienceResponse, error) {
	return s.userClient.PostExperienceRequest(ctx, in)
}
