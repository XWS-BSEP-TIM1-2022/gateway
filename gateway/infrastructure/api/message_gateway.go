package api

import (
	"errors"
	"fmt"
	"gateway/startup/config"
	messageService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/message"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

type MessageGatewayStruct struct {
	messageService.UnimplementedMessageServiceServer
	config        *config.Config
	messageClient messageService.MessageServiceClient
	userClient    userService.UserServiceClient
}

func NewMessageGateway(c *config.Config) *MessageGatewayStruct {
	return &MessageGatewayStruct{
		config:        c,
		messageClient: services.NewMessageClient(fmt.Sprintf("%s:%s", c.MessageServiceHost, c.MessageServicePort)),
		userClient:    services.NewUserClient(fmt.Sprintf("%s:%s", c.UserServiceHost, c.UserServicePort)),
	}
}

func (s *MessageGatewayStruct) GetAllNotifications(ctx context.Context, in *messageService.UserIdRequest) (*messageService.GetAllResponse, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &messageService.GetAllResponse{}, err
	}
	err = s.roleHavePermission(role, "notification_read")
	if err != nil {
		return &messageService.GetAllResponse{}, err
	}

	return s.messageClient.GetAllNotifications(ctx, in)
}

/*func (s *MessageGatewayStruct) CreateNotification(ctx context.Context, in *messageService.NewNotificationRequest) (*messageService.Notification, error) {

}*/

func (s *MessageGatewayStruct) isUserAuthenticated(ctx context.Context) (string, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return "", errors.New("unauthorized")
	}
	role, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{Token: jwt[0]})
	if err != nil {
		return "", errors.New("unauthorized")
	}

	return role.UserRole, nil
}

func (s *MessageGatewayStruct) roleHavePermission(role string, requiredPermission string) error {
	permissions := s.config.RolePermissions[role]
	if !contains(permissions, requiredPermission) {
		return errors.New("unauthorized")
	}

	return nil
}
