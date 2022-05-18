package api

import (
	"context"
	"errors"
	"fmt"
	"gateway/startup/config"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
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
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	_, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{Token: jwt[0]})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.userClient.UpdateRequest(ctx, in)
}

func (s *UserGatewayStruct) DeleteRequest(ctx context.Context, in *user.UserIdRequest) (*user.EmptyRequest, error) {
	return s.userClient.DeleteRequest(ctx, in)
}

func (s *UserGatewayStruct) ConfirmRegistration(ctx context.Context, in *user.ConfirmationRequest) (*user.ConfirmationResponse, error) {
	return s.userClient.ConfirmRegistration(ctx, in)
}

func (s *UserGatewayStruct) LoginRequest(ctx context.Context, in *user.CredentialsRequest) (*user.LoginResponse, error) {
	return s.userClient.LoginRequest(ctx, in)
}

func (s *UserGatewayStruct) GetQR2FA(ctx context.Context, in *user.UserIdRequest) (*user.TFAResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	_, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{Token: jwt[0]})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.userClient.GetQR2FA(ctx, in)
}

func (s *UserGatewayStruct) Enable2FA(ctx context.Context, in *user.TFARequest) (*user.EmptyRequest, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	_, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{Token: jwt[0]})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.userClient.Enable2FA(ctx, in)
}

func (s *UserGatewayStruct) Verify2FA(ctx context.Context, in *user.TFARequest) (*user.LoginResponse, error) {
	return s.userClient.Verify2FA(ctx, in)
}

func (s *UserGatewayStruct) Disable2FA(ctx context.Context, in *user.UserIdRequest) (*user.EmptyRequest, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	_, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{Token: jwt[0]})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.userClient.Disable2FA(ctx, in)
}

func (s *UserGatewayStruct) SearchUsersRequest(ctx context.Context, in *user.SearchRequest) (*user.UsersResponse, error) {
	return s.userClient.SearchUsersRequest(ctx, in)
}

func (s *UserGatewayStruct) IsUserAuthenticated(ctx context.Context, in *userService.AuthRequest) (*userService.AuthResponse, error) {
	return s.userClient.IsUserAuthenticated(ctx, in)
}

func (s *UserGatewayStruct) IsApiTokenValid(ctx context.Context, in *userService.AuthRequest) (*userService.UserIdRequest, error) {
	return s.userClient.IsApiTokenValid(ctx, in)
}

func (s *UserGatewayStruct) UpdatePasswordRequest(ctx context.Context, in *userService.NewPasswordRequest) (*user.GetResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	_, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{Token: jwt[0]})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.userClient.UpdatePasswordRequest(ctx, in)
}

func (s *UserGatewayStruct) GetAllUsersExperienceRequest(ctx context.Context, in *userService.ExperienceRequest) (*user.ExperienceResponse, error) {
	return s.userClient.GetAllUsersExperienceRequest(ctx, in)
}

func (s *UserGatewayStruct) PostExperienceRequest(ctx context.Context, in *user.NewExperienceRequest) (*user.NewExperienceResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	_, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{Token: jwt[0]})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.userClient.PostExperienceRequest(ctx, in)
}

func (s *UserGatewayStruct) DeleteExperienceRequest(ctx context.Context, in *user.DeleteUsersExperienceRequest) (*user.EmptyRequest, error) {
	return s.userClient.DeleteExperienceRequest(ctx, in)
}

func (s *UserGatewayStruct) AddUserSkill(ctx context.Context, in *user.NewSkillRequest) (*user.EmptyRequest, error) {
	return s.userClient.AddUserSkill(ctx, in)
}
func (s *UserGatewayStruct) AddUserInterest(ctx context.Context, in *user.NewInterestRequest) (*user.EmptyRequest, error) {
	return s.userClient.AddUserInterest(ctx, in)
}

func (s *UserGatewayStruct) RemoveInterest(ctx context.Context, in *user.RemoveInterestRequest) (*user.EmptyRequest, error) {
	return s.userClient.RemoveInterest(ctx, in)
}

func (s *UserGatewayStruct) RemoveSkill(ctx context.Context, in *user.RemoveSkillRequest) (*user.EmptyRequest, error) {
	return s.userClient.RemoveSkill(ctx, in)
}

func (s *UserGatewayStruct) ApiTokenRequest(ctx context.Context, in *user.UserIdRequest) (*user.ApiTokenResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	_, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{Token: jwt[0]})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.userClient.ApiTokenRequest(ctx, in)
}

func (s *UserGatewayStruct) ApiTokenCreateRequest(ctx context.Context, in *user.UserIdRequest) (*user.ApiTokenResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	_, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{Token: jwt[0]})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.userClient.ApiTokenCreateRequest(ctx, in)
}

func (s *UserGatewayStruct) ApiTokenRemoveRequest(ctx context.Context, in *user.UserIdRequest) (*user.EmptyRequest, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	_, err := s.userClient.IsUserAuthenticated(ctx, &userService.AuthRequest{Token: jwt[0]})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	return s.userClient.ApiTokenRemoveRequest(ctx, in)
}

func (s *UserGatewayStruct) CreatePasswordRecoveryRequest(ctx context.Context, in *user.UsernameRequest) (*user.EmptyRequest, error) {
	return s.userClient.CreatePasswordRecoveryRequest(ctx, in)
}

func (s *UserGatewayStruct) PasswordRecoveryRequest(ctx context.Context, in *user.NewPasswordRecoveryRequest) (*user.EmptyRequest, error) {
	return s.userClient.PasswordRecoveryRequest(ctx, in)
}
