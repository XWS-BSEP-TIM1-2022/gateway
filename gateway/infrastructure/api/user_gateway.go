package api

import (
	"context"
	"errors"
	"fmt"
	"gateway/startup/config"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/token"
	"google.golang.org/grpc/metadata"
	"regexp"
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
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.UsersResponse{}, err
	}
	err = s.roleHavePermission(role, "user_getAll")
	if err != nil {
		return &user.UsersResponse{}, err
	}

	return s.userClient.GetAllRequest(ctx, in)
}

func (s *UserGatewayStruct) PostRequest(ctx context.Context, in *user.UserRequest) (*user.GetResponse, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	return s.userClient.PostRequest(ctx, in)
}

func (s *UserGatewayStruct) PostAdminRequest(ctx context.Context, in *user.UserRequest) (*user.GetResponse, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	return s.userClient.PostAdminRequest(ctx, in)
}

func (s *UserGatewayStruct) UpdateRequest(ctx context.Context, in *user.UserRequest) (*user.GetResponse, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.GetResponse{}, err
	}
	err = s.roleHavePermission(role, "user_write")
	if err != nil {
		return &user.GetResponse{}, err
	}

	return s.userClient.UpdateRequest(ctx, in)
}

func (s *UserGatewayStruct) DeleteRequest(ctx context.Context, in *user.UserIdRequest) (*user.EmptyRequest, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "user_delete")
	if err != nil {
		return &user.EmptyRequest{}, err
	}

	return s.userClient.DeleteRequest(ctx, in)
}

func (s *UserGatewayStruct) ConfirmRegistration(ctx context.Context, in *user.ConfirmationRequest) (*user.ConfirmationResponse, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	return s.userClient.ConfirmRegistration(ctx, in)
}

func (s *UserGatewayStruct) LoginRequest(ctx context.Context, in *user.CredentialsRequest) (*user.LoginResponse, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	return s.userClient.LoginRequest(ctx, in)
}

func (s *UserGatewayStruct) GetQR2FA(ctx context.Context, in *user.UserIdRequest) (*user.TFAResponse, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.TFAResponse{}, err
	}
	err = s.roleHavePermission(role, "user_read")
	if err != nil {
		return &user.TFAResponse{}, err
	}

	return s.userClient.GetQR2FA(ctx, in)
}

func (s *UserGatewayStruct) Enable2FA(ctx context.Context, in *user.TFARequest) (*user.EmptyRequest, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "user_write")
	if err != nil {
		return &user.EmptyRequest{}, err
	}

	return s.userClient.Enable2FA(ctx, in)
}

func (s *UserGatewayStruct) Verify2FA(ctx context.Context, in *user.TFARequest) (*user.LoginResponse, error) {
	return s.userClient.Verify2FA(ctx, in)
}

func (s *UserGatewayStruct) Disable2FA(ctx context.Context, in *user.UserIdRequest) (*user.EmptyRequest, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "user_write")
	if err != nil {
		return &user.EmptyRequest{}, err
	}

	return s.userClient.Disable2FA(ctx, in)
}

func (s *UserGatewayStruct) SearchUsersRequest(ctx context.Context, in *user.SearchRequest) (*user.UsersResponse, error) {
	in.UserId = getUserIdFromJwt(ctx)

	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	return s.userClient.SearchUsersRequest(ctx, in)
}

func (s *UserGatewayStruct) IsUserAuthenticated(ctx context.Context, in *userService.AuthRequest) (*userService.AuthResponse, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	return s.userClient.IsUserAuthenticated(ctx, in)
}

func (s *UserGatewayStruct) IsApiTokenValid(ctx context.Context, in *userService.AuthRequest) (*userService.UserIdRequest, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	return s.userClient.IsApiTokenValid(ctx, in)
}

func (s *UserGatewayStruct) UpdatePasswordRequest(ctx context.Context, in *userService.NewPasswordRequest) (*user.GetResponse, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.GetResponse{}, err
	}
	err = s.roleHavePermission(role, "user_write")
	if err != nil {
		return &user.GetResponse{}, err
	}

	return s.userClient.UpdatePasswordRequest(ctx, in)
}

func (s *UserGatewayStruct) GetAllUsersExperienceRequest(ctx context.Context, in *userService.ExperienceRequest) (*user.ExperienceResponse, error) {
	return s.userClient.GetAllUsersExperienceRequest(ctx, in)
}

func (s *UserGatewayStruct) PostExperienceRequest(ctx context.Context, in *user.NewExperienceRequest) (*user.NewExperienceResponse, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.NewExperienceResponse{}, err
	}
	err = s.roleHavePermission(role, "user_write")
	if err != nil {
		return &user.NewExperienceResponse{}, err
	}

	return s.userClient.PostExperienceRequest(ctx, in)
}

func (s *UserGatewayStruct) DeleteExperienceRequest(ctx context.Context, in *user.DeleteUsersExperienceRequest) (*user.EmptyRequest, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "user_write")
	if err != nil {
		return &user.EmptyRequest{}, err
	}

	return s.userClient.DeleteExperienceRequest(ctx, in)
}

func (s *UserGatewayStruct) AddUserSkill(ctx context.Context, in *user.NewSkillRequest) (*user.EmptyRequest, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "user_write")
	if err != nil {
		return &user.EmptyRequest{}, err
	}

	return s.userClient.AddUserSkill(ctx, in)
}
func (s *UserGatewayStruct) AddUserInterest(ctx context.Context, in *user.NewInterestRequest) (*user.EmptyRequest, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "user_write")
	if err != nil {
		return &user.EmptyRequest{}, err
	}

	return s.userClient.AddUserInterest(ctx, in)
}

func (s *UserGatewayStruct) RemoveInterest(ctx context.Context, in *user.RemoveInterestRequest) (*user.EmptyRequest, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "user_write")
	if err != nil {
		return &user.EmptyRequest{}, err
	}

	return s.userClient.RemoveInterest(ctx, in)
}

func (s *UserGatewayStruct) RemoveSkill(ctx context.Context, in *user.RemoveSkillRequest) (*user.EmptyRequest, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "user_write")
	if err != nil {
		return &user.EmptyRequest{}, err
	}

	return s.userClient.RemoveSkill(ctx, in)
}

func (s *UserGatewayStruct) ApiTokenRequest(ctx context.Context, in *user.UserIdRequest) (*user.ApiTokenResponse, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.ApiTokenResponse{}, err
	}
	err = s.roleHavePermission(role, "user_read")
	if err != nil {
		return &user.ApiTokenResponse{}, err
	}

	return s.userClient.ApiTokenRequest(ctx, in)
}

func (s *UserGatewayStruct) ApiTokenCreateRequest(ctx context.Context, in *user.UserIdRequest) (*user.ApiTokenResponse, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.ApiTokenResponse{}, err
	}
	err = s.roleHavePermission(role, "user_write")
	if err != nil {
		return &user.ApiTokenResponse{}, err
	}

	return s.userClient.ApiTokenCreateRequest(ctx, in)
}

func (s *UserGatewayStruct) ApiTokenRemoveRequest(ctx context.Context, in *user.UserIdRequest) (*user.EmptyRequest, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "user_write")
	if err != nil {
		return &user.EmptyRequest{}, err
	}

	return s.userClient.ApiTokenRemoveRequest(ctx, in)
}

func (s *UserGatewayStruct) CreatePasswordRecoveryRequest(ctx context.Context, in *user.UsernameRequest) (*user.EmptyRequest, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	return s.userClient.CreatePasswordRecoveryRequest(ctx, in)
}

func (s *UserGatewayStruct) PasswordRecoveryRequest(ctx context.Context, in *user.NewPasswordRecoveryRequest) (*user.EmptyRequest, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	return s.userClient.PasswordRecoveryRequest(ctx, in)
}

func (s *UserGatewayStruct) PasswordlessLoginStart(ctx context.Context, in *user.UsernameRequest) (*user.EmptyRequest, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	return s.userClient.PasswordlessLoginStart(ctx, in)
}

func (s *UserGatewayStruct) PasswordlessLogin(ctx context.Context, in *user.PasswordlessLoginRequest) (*user.LoginResponse, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	return s.userClient.PasswordlessLogin(ctx, in)
}

func (s *UserGatewayStruct) ChangeProfilePrivacy(ctx context.Context, in *user.UserIdRequest) (*user.EmptyRequest, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &user.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "user_write")
	if err != nil {
		return &user.EmptyRequest{}, err
	}

	return s.userClient.ChangeProfilePrivacy(ctx, in)
}

func (s *UserGatewayStruct) isUserAuthenticated(ctx context.Context) (string, error) {
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

func (s *UserGatewayStruct) roleHavePermission(role string, requiredPermission string) error {
	permissions := s.config.RolePermissions[role]
	if !contains(permissions, requiredPermission) {
		return errors.New("unauthorized")
	}

	return nil
}

func checkValue(value string) error {
	match, _ := regexp.MatchString("<[^\\w<>]*(?:[^<>\"'\\s]*:)?[^\\w<>]*(?:\\W*s\\W*c\\W*r\\W*i\\W*p\\W*t|\\W*f\\W*o\\W*r\\W*m|\\W*s\\W*t\\W*y\\W*l\\W*e|\\W*s\\W*v\\W*g|\\W*m\\W*a\\W*r\\W*q\\W*u\\W*e\\W*e|(?:\\W*l\\W*i\\W*n\\W*k|\\W*o\\W*b\\W*j\\W*e\\W*c\\W*t|\\W*e\\W*m\\W*b\\W*e\\W*d|\\W*a\\W*p\\W*p\\W*l\\W*e\\W*t|\\W*p\\W*a\\W*r\\W*a\\W*m|\\W*i?\\W*f\\W*r\\W*a\\W*m\\W*e|\\W*b\\W*a\\W*s\\W*e|\\W*b\\W*o\\W*d\\W*y|\\W*m\\W*e\\W*t\\W*a|\\W*i\\W*m\\W*a?\\W*g\\W*e?|\\W*v\\W*i\\W*d\\W*e\\W*o|\\W*a\\W*u\\W*d\\W*i\\W*o|\\W*b\\W*i\\W*n\\W*d\\W*i\\W*n\\W*g\\W*s|\\W*s\\W*e\\W*t|\\W*i\\W*s\\W*i\\W*n\\W*d\\W*e\\W*x|\\W*a\\W*n\\W*i\\W*m\\W*a\\W*t\\W*e)[^>\\w])|(?:<\\w[\\s\\S]*[\\s\\0\\/]|['\"])(?:formaction|style|background|src|lowsrc|ping|on(?:d(?:e(?:vice(?:(?:orienta|mo)tion|proximity|found|light)|livery(?:success|error)|activate)|r(?:ag(?:e(?:n(?:ter|d)|xit)|(?:gestur|leav)e|start|drop|over)?|op)|i(?:s(?:c(?:hargingtimechange|onnect(?:ing|ed))|abled)|aling)|ata(?:setc(?:omplete|hanged)|(?:availabl|chang)e|error)|urationchange|ownloading|blclick)|Moz(?:M(?:agnifyGesture(?:Update|Start)?|ouse(?:PixelScroll|Hittest))|S(?:wipeGesture(?:Update|Start|End)?|crolledAreaChanged)|(?:(?:Press)?TapGestur|BeforeResiz)e|EdgeUI(?:C(?:omplet|ancel)|Start)ed|RotateGesture(?:Update|Start)?|A(?:udioAvailable|fterPaint))|c(?:o(?:m(?:p(?:osition(?:update|start|end)|lete)|mand(?:update)?)|n(?:t(?:rolselect|extmenu)|nect(?:ing|ed))|py)|a(?:(?:llschang|ch)ed|nplay(?:through)?|rdstatechange)|h(?:(?:arging(?:time)?ch)?ange|ecking)|(?:fstate|ell)change|u(?:echange|t)|l(?:ick|ose))|m(?:o(?:z(?:pointerlock(?:change|error)|(?:orientation|time)change|fullscreen(?:change|error)|network(?:down|up)load)|use(?:(?:lea|mo)ve|o(?:ver|ut)|enter|wheel|down|up)|ve(?:start|end)?)|essage|ark)|s(?:t(?:a(?:t(?:uschanged|echange)|lled|rt)|k(?:sessione|comma)nd|op)|e(?:ek(?:complete|ing|ed)|(?:lec(?:tstar)?)?t|n(?:ding|t))|u(?:ccess|spend|bmit)|peech(?:start|end)|ound(?:start|end)|croll|how)|b(?:e(?:for(?:e(?:(?:scriptexecu|activa)te|u(?:nload|pdate)|p(?:aste|rint)|c(?:opy|ut)|editfocus)|deactivate)|gin(?:Event)?)|oun(?:dary|ce)|l(?:ocked|ur)|roadcast|usy)|a(?:n(?:imation(?:iteration|start|end)|tennastatechange)|fter(?:(?:scriptexecu|upda)te|print)|udio(?:process|start|end)|d(?:apteradded|dtrack)|ctivate|lerting|bort)|DOM(?:Node(?:Inserted(?:IntoDocument)?|Removed(?:FromDocument)?)|(?:CharacterData|Subtree)Modified|A(?:ttrModified|ctivate)|Focus(?:Out|In)|MouseScroll)|r(?:e(?:s(?:u(?:m(?:ing|e)|lt)|ize|et)|adystatechange|pea(?:tEven)?t|movetrack|trieving|ceived)|ow(?:s(?:inserted|delete)|e(?:nter|xit))|atechange)|p(?:op(?:up(?:hid(?:den|ing)|show(?:ing|n))|state)|a(?:ge(?:hide|show)|(?:st|us)e|int)|ro(?:pertychange|gress)|lay(?:ing)?)|t(?:ouch(?:(?:lea|mo)ve|en(?:ter|d)|cancel|start)|ime(?:update|out)|ransitionend|ext)|u(?:s(?:erproximity|sdreceived)|p(?:gradeneeded|dateready)|n(?:derflow|load))|f(?:o(?:rm(?:change|input)|cus(?:out|in)?)|i(?:lterchange|nish)|ailed)|l(?:o(?:ad(?:e(?:d(?:meta)?data|nd)|start)?|secapture)|evelchange|y)|g(?:amepad(?:(?:dis)?connected|button(?:down|up)|axismove)|et)|e(?:n(?:d(?:Event|ed)?|abled|ter)|rror(?:update)?|mptied|xit)|i(?:cc(?:cardlockerror|infochange)|n(?:coming|valid|put))|o(?:(?:(?:ff|n)lin|bsolet)e|verflow(?:changed)?|pen)|SVG(?:(?:Unl|L)oad|Resize|Scroll|Abort|Error|Zoom)|h(?:e(?:adphoneschange|l[dp])|ashchange|olding)|v(?:o(?:lum|ic)e|ersion)change|w(?:a(?:it|rn)ing|heel)|key(?:press|down|up)|(?:AppComman|Loa)d|no(?:update|match)|Request|zoom))[\\s\\0]*=\n", value)
	if match {
		return errors.New("forbidden stuff in input")
	}
	return nil
}

func getUserIdFromJwt(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	retVal := ""
	if jwt != nil {
		var err error
		retVal, err = token.NewJwtManagerDislinkt(0).GetUserIdFromToken(jwt[0])
		if err != nil {
			retVal = ""
		}
	}
	return retVal
}
