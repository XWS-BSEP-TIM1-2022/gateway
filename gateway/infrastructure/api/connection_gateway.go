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
	Log.Info("Creating new connection for user with id: " + in.Connection.UserId + " with user with id: " + in.Connection.ConnectedUserId)
	err := checkValue(in.String())
	if err != nil {
		Log.Warn("Input possibly contains malicious data")
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("User is not authenticated")
		return &connectionService.UserConnectionResponse{}, err
	}
	err = s.roleHavePermission(role, "connection_write")
	if err != nil {
		Log.Warn("Current user role dont have valid permission")
		return &connectionService.UserConnectionResponse{}, err
	}

	return s.connectionClient.NewUserConnection(ctx, in)
}

func (s *ConnectionGatewayStruct) ApproveConnection(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	Log.Info("Approving connection between users with ids: " + in.Connection.UserId + "---" + in.Connection.ConnectedUserId)
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.UserConnectionResponse{}, err
	}
	err = s.roleHavePermission(role, "connection_write")
	if err != nil {
		return &connectionService.UserConnectionResponse{}, err
	}

	return s.connectionClient.ApproveConnection(ctx, in)
}

func (s *ConnectionGatewayStruct) GetConnection(ctx context.Context, in *connectionService.Connection) (*connectionService.Connection, error) {
	Log.Info("Getting connection between users with ids: " + in.UserId + "---" + in.ConnectedUserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.Connection{}, err
	}
	err = s.roleHavePermission(role, "connection_read")
	if err != nil {
		return &connectionService.Connection{}, err
	}

	return s.connectionClient.GetConnection(ctx, in)
}

func (s *ConnectionGatewayStruct) ApproveAllConnection(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.EmptyRequest, error) {
	Log.Info("Approving all connections for user with id: " + in.UserId)
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "connection_write")
	if err != nil {
		return &connectionService.EmptyRequest{}, err
	}

	return s.connectionClient.ApproveAllConnection(ctx, in)
}

func (s *ConnectionGatewayStruct) RejectConnection(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	Log.Info("Getting connection between users with ids: " + in.Connection.UserId + "---" + in.Connection.ConnectedUserId)
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.UserConnectionResponse{}, err
	}
	err = s.roleHavePermission(role, "connection_write")
	if err != nil {
		return &connectionService.UserConnectionResponse{}, err
	}

	return s.connectionClient.RejectConnection(ctx, in)
}

func (s *ConnectionGatewayStruct) DeleteConnection(ctx context.Context, in *connectionService.Connection) (*connectionService.UserConnectionResponse, error) {
	Log.Info("Delete connection between users with ids: " + in.UserId + "---" + in.ConnectedUserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.UserConnectionResponse{}, err
	}
	err = s.roleHavePermission(role, "connection_delete")
	if err != nil {
		return &connectionService.UserConnectionResponse{}, err
	}

	return s.connectionClient.DeleteConnection(ctx, in)
}

func (s *ConnectionGatewayStruct) GetAllConnections(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	Log.Info("Getting all connections")
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.AllConnectionResponse{}, err
	}
	err = s.roleHavePermission(role, "connection_read")
	if err != nil {
		return &connectionService.AllConnectionResponse{}, err
	}

	return s.connectionClient.GetAllConnections(ctx, in)
}

func (s *ConnectionGatewayStruct) GetFollowings(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	Log.Info("Getting all following for user with id: " + in.UserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.AllConnectionResponse{}, err
	}
	err = s.roleHavePermission(role, "connection_read")
	if err != nil {
		return &connectionService.AllConnectionResponse{}, err
	}

	return s.connectionClient.GetFollowings(ctx, in)
}

func (s *ConnectionGatewayStruct) GetFollowers(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	Log.Info("Getting all followers for user with id: " + in.UserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.AllConnectionResponse{}, err
	}
	err = s.roleHavePermission(role, "connection_read")
	if err != nil {
		return &connectionService.AllConnectionResponse{}, err
	}

	return s.connectionClient.GetFollowers(ctx, in)
}

func (s *ConnectionGatewayStruct) GetAllRequestConnectionsByUserId(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	Log.Info("Getting connection request for user with id: " + in.UserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.AllConnectionResponse{}, err
	}
	err = s.roleHavePermission(role, "connection_read")
	if err != nil {
		return &connectionService.AllConnectionResponse{}, err
	}

	return s.connectionClient.GetAllRequestConnectionsByUserId(ctx, in)
}

func (s *ConnectionGatewayStruct) GetAllPendingConnectionsByUserId(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	Log.Info("Getting all pending request for user with id: " + in.UserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.AllConnectionResponse{}, err
	}
	err = s.roleHavePermission(role, "connection_read")
	if err != nil {
		return &connectionService.AllConnectionResponse{}, err
	}

	return s.connectionClient.GetAllPendingConnectionsByUserId(ctx, in)
}

func (s *ConnectionGatewayStruct) BlockUser(ctx context.Context, in *connectionService.BlockUserRequest) (*connectionService.EmptyRequest, error) {
	Log.Info("User with id: " + in.Block.UserId + "block user with id: " + in.Block.BlockUserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "block_write")
	if err != nil {
		return &connectionService.EmptyRequest{}, err
	}

	return s.connectionClient.BlockUser(ctx, in)
}

func (s *ConnectionGatewayStruct) UnblockUser(ctx context.Context, in *connectionService.BlockUserRequest) (*connectionService.EmptyRequest, error) {
	Log.Info("User with id: " + in.Block.UserId + "unblock user with id: " + in.Block.BlockUserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "block_write")
	if err != nil {
		return &connectionService.EmptyRequest{}, err
	}

	return s.connectionClient.UnblockUser(ctx, in)
}

func (s *ConnectionGatewayStruct) IsBlocked(ctx context.Context, in *connectionService.Block) (*connectionService.IsBlockedResponse, error) {
	Log.Info("User with id: " + in.UserId + "check is user with id: " + in.BlockUserId + " blocked")
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.IsBlockedResponse{}, err
	}
	err = s.roleHavePermission(role, "block_read")
	if err != nil {
		return &connectionService.IsBlockedResponse{}, err
	}

	return s.connectionClient.IsBlocked(ctx, in)
}

func (s *ConnectionGatewayStruct) IsBlockedAny(ctx context.Context, in *connectionService.Block) (*connectionService.IsBlockedResponse, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.IsBlockedResponse{}, err
	}
	err = s.roleHavePermission(role, "block_read")
	if err != nil {
		return &connectionService.IsBlockedResponse{}, err
	}
	return s.connectionClient.IsBlockedAny(ctx, in)
}

func (s *ConnectionGatewayStruct) Blocked(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.BlockedResponse, error) {
	Log.Info("Check is user with id:" + in.UserId + " blocked")
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.BlockedResponse{}, err
	}
	err = s.roleHavePermission(role, "block_read")
	if err != nil {
		return &connectionService.BlockedResponse{}, err
	}
	return s.connectionClient.Blocked(ctx, in)
}

func (s *ConnectionGatewayStruct) BlockedBy(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.BlockedResponse, error) {
	Log.Info("Checking is current user blocked by user with id:" + in.UserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.BlockedResponse{}, err
	}
	err = s.roleHavePermission(role, "block_read")
	if err != nil {
		return &connectionService.BlockedResponse{}, err
	}
	return s.connectionClient.BlockedBy(ctx, in)
}

func (s *ConnectionGatewayStruct) BlockedAny(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.BlockedResponse, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.BlockedResponse{}, err
	}
	err = s.roleHavePermission(role, "block_read")
	if err != nil {
		return &connectionService.BlockedResponse{}, err
	}
	return s.connectionClient.BlockedAny(ctx, in)
}

func (s *ConnectionGatewayStruct) ChangeMessageNotification(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	Log.Info("Changing message notification with user id: " + in.Connection.UserId + " for user with id:" + in.Connection.ConnectedUserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.UserConnectionResponse{}, err
	}
	err = s.roleHavePermission(role, "connection_write")
	if err != nil {
		return &connectionService.UserConnectionResponse{}, err
	}

	return s.connectionClient.ChangeMessageNotification(ctx, in)
}

func (s *ConnectionGatewayStruct) ChangePostNotification(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	Log.Info("Changing post notification with user id: " + in.Connection.UserId + " for user with id:" + in.Connection.ConnectedUserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.UserConnectionResponse{}, err
	}
	err = s.roleHavePermission(role, "connection_write")
	if err != nil {
		return &connectionService.UserConnectionResponse{}, err
	}

	return s.connectionClient.ChangePostNotification(ctx, in)
}

func (s *ConnectionGatewayStruct) ChangeCommentNotification(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	Log.Info("Changing comment notification with user id: " + in.Connection.UserId + " for user with id:" + in.Connection.ConnectedUserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &connectionService.UserConnectionResponse{}, err
	}
	err = s.roleHavePermission(role, "connection_write")
	if err != nil {
		return &connectionService.UserConnectionResponse{}, err
	}

	return s.connectionClient.ChangeCommentNotification(ctx, in)
}

func (s *ConnectionGatewayStruct) isUserAuthenticated(ctx context.Context) (string, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		Log.Warn("User dont have jwt in request")
		return "", errors.New("unauthorized")
	}
	role, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{Token: jwt[0]})
	if err != nil {
		Log.Warn("User is not authenticated")
		return "", errors.New("unauthorized")
	}

	return role.UserRole, nil
}

func (s *ConnectionGatewayStruct) roleHavePermission(role string, requiredPermission string) error {
	permissions := s.config.RolePermissions[role]
	if !contains(permissions, requiredPermission) {
		Log.Warn("User doesn't have permission to get requests")
		return errors.New("unauthorized")
	}

	return nil
}
