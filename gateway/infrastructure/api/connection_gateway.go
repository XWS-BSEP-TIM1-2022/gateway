package api

import (
	"context"
	"errors"
	"fmt"
	"gateway/startup/config"
	connectionService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/connection"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
	"google.golang.org/grpc/metadata"
)

type ConnectionGatewayStruct struct {
	connectionService.UnimplementedConnectionServiceServer
	config           *config.Config
	connectionClient connectionService.ConnectionServiceClient
	userClient       userService.UserServiceClient
}

func NewConnectionGateway(c *config.Config) *ConnectionGatewayStruct {
	return &ConnectionGatewayStruct{
		config:           c,
		connectionClient: services.NewConnectionClient(fmt.Sprintf("%s:%s", c.ConnectionServiceHost, c.ConnectionServicePort)),
		userClient:       services.NewUserClient(fmt.Sprintf("%s:%s", c.UserServiceHost, c.UserServicePort)),
	}
}

func (s *ConnectionGatewayStruct) NewUserConnection(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	role, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{JwtToken: jwt[0]})
	if role.UserRole != "USER" || err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.connectionClient.NewUserConnection(ctx, in)
}

func (s *ConnectionGatewayStruct) ApproveConnection(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	role, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{JwtToken: jwt[0]})
	if role.UserRole != "USER" || err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.connectionClient.ApproveConnection(ctx, in)
}

func (s *ConnectionGatewayStruct) RejectConnection(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	role, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{JwtToken: jwt[0]})
	if role.UserRole != "USER" || err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.connectionClient.RejectConnection(ctx, in)
}

func (s *ConnectionGatewayStruct) DeleteConnection(ctx context.Context, in *connectionService.Connection) (*connectionService.UserConnectionResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	role, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{JwtToken: jwt[0]})
	if role.UserRole != "USER" || err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.connectionClient.DeleteConnection(ctx, in)
}

func (s *ConnectionGatewayStruct) GetAllConnections(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	role, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{JwtToken: jwt[0]})
	if role.UserRole != "USER" || err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.connectionClient.GetAllConnections(ctx, in)
}

func (s *ConnectionGatewayStruct) GetFollowings(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	role, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{JwtToken: jwt[0]})
	if role.UserRole != "USER" || err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.connectionClient.GetFollowings(ctx, in)
}

func (s *ConnectionGatewayStruct) GetFollowers(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	role, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{JwtToken: jwt[0]})
	if role.UserRole != "USER" || err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.connectionClient.GetFollowers(ctx, in)
}

func (s *ConnectionGatewayStruct) GetAllRequestConnectionsByUserId(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	role, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{JwtToken: jwt[0]})
	if role.UserRole != "USER" || err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.connectionClient.GetAllRequestConnectionsByUserId(ctx, in)
}

func (s *ConnectionGatewayStruct) GetAllPendingConnectionsByUserId(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	role, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{JwtToken: jwt[0]})
	if role.UserRole != "USER" || err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.connectionClient.GetAllPendingConnectionsByUserId(ctx, in)
}
