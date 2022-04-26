package api

import (
	"context"
	"fmt"
	"gateway/startup/config"
	connectionService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/connection"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
)

type ConnectionGatewayStruct struct {
	connectionService.UnimplementedConnectionServiceServer
	config           *config.Config
	connectionClient connectionService.ConnectionServiceClient
}

func NewConnectionGateway(c *config.Config) *ConnectionGatewayStruct {
	return &ConnectionGatewayStruct{
		config:           c,
		connectionClient: services.NewConnectionClient(fmt.Sprintf("%s:%s", c.ConnectionServiceHost, c.ConnectionServicePort)),
	}
}

func (s *ConnectionGatewayStruct) NewUserConnection(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	return s.connectionClient.NewUserConnection(ctx, in)
}

func (s *ConnectionGatewayStruct) ApproveConnection(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	return s.connectionClient.ApproveConnection(ctx, in)
}
