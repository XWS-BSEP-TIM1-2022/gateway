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
	Log.Info("Getting all notifications for user with id: " + in.UserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("Unauthenticated request for user with id: " + in.UserId)
		return &messageService.GetAllResponse{}, err
	}
	err = s.roleHavePermission(role, "notification_read")
	if err != nil {
		Log.Warn("User with id: " + in.UserId + " doesn't have permission to get requests")
		return &messageService.GetAllResponse{}, err
	}

	return s.messageClient.GetAllNotifications(ctx, in)
}

/*func (s *MessageGatewayStruct) CreateNotification(ctx context.Context, in *messageService.NewNotificationRequest) (*messageService.Notification, error) {

}*/

func (s *MessageGatewayStruct) GetAllMessagesForUser(ctx context.Context, in *messageService.ChatIdRequest) (*messageService.GetAllMessagesResponse, error) {
	Log.Info("Getting all messages for chat with id: " + in.ChatId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("Unauthenticated request for chat with id: " + in.ChatId)
		return &messageService.GetAllMessagesResponse{}, err
	}
	err = s.roleHavePermission(role, "message_read")
	if err != nil {
		Log.Warn("User doesn't have permission to get requests for chat with id:" + in.ChatId)
		return &messageService.GetAllMessagesResponse{}, err
	}

	return s.messageClient.GetAllMessagesForUser(ctx, in)
}

func (s *MessageGatewayStruct) CreateMessage(ctx context.Context, in *messageService.NewMessageRequest) (*messageService.GetMessageResponse, error) {
	Log.Info("Creating new message")
	err := checkValue(in.Message.Message)
	if err != nil {
		Log.Warn("Input possibly contains malicious data")
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("User is not authenticated")
		return &messageService.GetMessageResponse{}, err
	}
	err = s.roleHavePermission(role, "message_write")
	if err != nil {
		Log.Warn("User doesn't have permission to get requests")
		return &messageService.GetMessageResponse{}, err
	}

	return s.messageClient.CreateMessage(ctx, in)
}

func (s *MessageGatewayStruct) GetAllChatsForUser(ctx context.Context, in *messageService.UserIdRequest) (*messageService.GetAllChatsResponse, error) {
	Log.Info("Getting all chats for user with id: " + in.UserId)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("Unauthenticated request for chat with id: " + in.UserId)
		return &messageService.GetAllChatsResponse{}, err
	}
	err = s.roleHavePermission(role, "chat_read")
	if err != nil {
		Log.Warn("User with id: " + in.UserId + " doesn't have permission to get requests")
		return &messageService.GetAllChatsResponse{}, err
	}

	return s.messageClient.GetAllChatsForUser(ctx, in)
}

func (s *MessageGatewayStruct) CreateChat(ctx context.Context, in *messageService.NewChatRequest) (*messageService.GetChatResponse, error) {
	Log.Info("Creating new chat")
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("User is not authenticated")
		return &messageService.GetChatResponse{}, err
	}
	err = s.roleHavePermission(role, "chat_write")
	if err != nil {
		Log.Warn("User doesn't have permission to get requests")
		return &messageService.GetChatResponse{}, err
	}

	return s.messageClient.CreateChat(ctx, in)
}

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
