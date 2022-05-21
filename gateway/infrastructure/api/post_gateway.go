package api

import (
	"context"
	"errors"
	"fmt"
	"gateway/startup/config"
	postService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/post"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
	"google.golang.org/grpc/metadata"
)

type PostGatewayStruct struct {
	postService.UnimplementedPostServiceServer
	config     *config.Config
	postClient postService.PostServiceClient
	userClient userService.UserServiceClient
}

func NewPostGateway(c *config.Config) *PostGatewayStruct {
	return &PostGatewayStruct{
		config:     c,
		postClient: services.NewPostClient(fmt.Sprintf("%s:%s", c.PostServiceHost, c.PostServicePort)),
		userClient: services.NewUserClient(fmt.Sprintf("%s:%s", c.UserServiceHost, c.UserServicePort)),
	}
}

func (s *PostGatewayStruct) GetRequest(ctx context.Context, in *postService.PostIdRequest) (*postService.PostResponse, error) {
	return s.postClient.GetRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllRequest(ctx context.Context, in *postService.EmptyRequest) (*postService.PostsResponse, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &postService.PostsResponse{}, err
	}
	err = s.roleHavePermission(role, "post_getAll")
	if err != nil {
		return &postService.PostsResponse{}, err
	}

	return s.postClient.GetAllRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllFromUserRequest(ctx context.Context, in *postService.UserPostsRequest) (*postService.PostsResponse, error) {
	return s.postClient.GetAllFromUserRequest(ctx, in)
}

func (s *PostGatewayStruct) CreateRequest(ctx context.Context, in *postService.PostRequest) (*postService.PostResponse, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &postService.PostResponse{}, err
	}
	err = s.roleHavePermission(role, "post_write")
	if err != nil {
		return &postService.PostResponse{}, err
	}

	return s.postClient.CreateRequest(ctx, in)
}

func (s *PostGatewayStruct) DeleteRequest(ctx context.Context, in *postService.PostIdRequest) (*postService.EmptyRequest, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &postService.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "post_delete")
	if err != nil {
		return &postService.EmptyRequest{}, err
	}

	return s.postClient.DeleteRequest(ctx, in)
}

func (s *PostGatewayStruct) GetCommentRequest(ctx context.Context, in *postService.CommentIdRequest) (*postService.CommentResponse, error) {
	return s.postClient.GetCommentRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllCommentsRequest(ctx context.Context, in *postService.EmptyRequest) (*postService.CommentsResponse, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &postService.CommentsResponse{}, err
	}
	err = s.roleHavePermission(role, "post_getAll")
	if err != nil {
		return &postService.CommentsResponse{}, err
	}

	return s.postClient.GetAllCommentsRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllCommentsFromPostRequest(ctx context.Context, in *postService.PostCommentsRequest) (*postService.CommentsResponse, error) {
	return s.postClient.GetAllCommentsFromPostRequest(ctx, in)
}

func (s *PostGatewayStruct) CreateCommentRequest(ctx context.Context, in *postService.CommentRequest) (*postService.CommentResponse, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &postService.CommentResponse{}, err
	}
	err = s.roleHavePermission(role, "post_write")
	if err != nil {
		return &postService.CommentResponse{}, err
	}

	return s.postClient.CreateCommentRequest(ctx, in)
}

func (s *PostGatewayStruct) DeleteCommentRequest(ctx context.Context, in *postService.CommentIdRequest) (*postService.EmptyRequest, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return &postService.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "post_delete")
	if err != nil {
		return &postService.EmptyRequest{}, err
	}
	return s.postClient.DeleteCommentRequest(ctx, in)
}

func (s *PostGatewayStruct) GetReactionRequest(ctx context.Context, in *postService.ReactionIdRequest) (*postService.ReactionResponse, error) {
	return s.postClient.GetReactionRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllReactionsRequest(ctx context.Context, in *postService.EmptyRequest) (*postService.ReactionsResponse, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return nil, err
	}
	err = s.roleHavePermission(role, "post_getAll")
	if err != nil {
		return nil, err
	}

	return s.postClient.GetAllReactionsRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllReactionsFromPostRequest(ctx context.Context, in *postService.PostReactionRequest) (*postService.ReactionsResponse, error) {
	return s.postClient.GetAllReactionsFromPostRequest(ctx, in)
}

func (s *PostGatewayStruct) CreateReactionRequest(ctx context.Context, in *postService.ReactionRequest) (*postService.ReactionResponse, error) {
	err := checkValue(in.String())
	if err != nil {
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return nil, err
	}
	err = s.roleHavePermission(role, "post_write")
	if err != nil {
		return nil, err
	}

	return s.postClient.CreateReactionRequest(ctx, in)
}

func (s *PostGatewayStruct) DeleteReactionRequest(ctx context.Context, in *postService.ReactionIdRequest) (*postService.EmptyRequest, error) {
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		return nil, err
	}
	err = s.roleHavePermission(role, "post_delete")
	if err != nil {
		return nil, err
	}

	return s.postClient.DeleteReactionRequest(ctx, in)
}

func (s *PostGatewayStruct) isUserAuthenticated(ctx context.Context) (string, error) {
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

func (s *PostGatewayStruct) roleHavePermission(role string, requiredPermission string) error {
	permissions := s.config.RolePermissions[role]
	if !contains(permissions, requiredPermission) {
		return errors.New("unauthorized")
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
