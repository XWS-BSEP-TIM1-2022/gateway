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
	Log.Info("Getting post by id")
	in.LoggedUserId = getUserIdFromJwt(ctx)
	return s.postClient.GetRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllRequest(ctx context.Context, in *postService.EmptyRequest) (*postService.PostsResponse, error) {
	Log.Info("Getting all posts")
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("User is not authenticated")
		return &postService.PostsResponse{}, err
	}
	err = s.roleHavePermission(role, "post_getAll")
	if err != nil {
		Log.Warn("User doesn't have permission to get requests")
		return &postService.PostsResponse{}, err
	}

	return s.postClient.GetAllRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllFromUserRequest(ctx context.Context, in *postService.UserPostsRequest) (*postService.PostsResponse, error) {
	Log.Info("Getting all users posts")
	in.LoggedUserId = getUserIdFromJwt(ctx)
	return s.postClient.GetAllFromUserRequest(ctx, in)
}

func (s *PostGatewayStruct) CreateRequest(ctx context.Context, in *postService.PostRequest) (*postService.PostResponse, error) {
	Log.Info("Creating new post")
	err := checkValue(in.String())
	if err != nil {
		Log.Warn("Input possibly contains malicious data")
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("User is not authenticated")
		return &postService.PostResponse{}, err
	}
	err = s.roleHavePermission(role, "post_write")
	if err != nil {
		Log.Warn("User doesn't have permission to get requests")
		return &postService.PostResponse{}, err
	}

	return s.postClient.CreateRequest(ctx, in)
}

func (s *PostGatewayStruct) DeleteRequest(ctx context.Context, in *postService.PostIdRequest) (*postService.EmptyRequest, error) {
	Log.Info("Deleting post")
	err := checkValue(in.String())
	if err != nil {
		Log.Warn("Input possibly contains malicious data")
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("User is not authenticated")
		return &postService.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "post_delete")
	if err != nil {
		Log.Warn("User doesn't have permission to get requests")
		return &postService.EmptyRequest{}, err
	}

	in.LoggedUserId = getUserIdFromJwt(ctx)

	return s.postClient.DeleteRequest(ctx, in)
}

func (s *PostGatewayStruct) GetCommentRequest(ctx context.Context, in *postService.CommentIdRequest) (*postService.CommentResponse, error) {
	Log.Info("User with id: " + in.LoggedUserId + " getting comment with id: " + in.Id)
	in.LoggedUserId = getUserIdFromJwt(ctx)
	return s.postClient.GetCommentRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllCommentsRequest(ctx context.Context, in *postService.EmptyRequest) (*postService.CommentsResponse, error) {
	Log.Info("Getting all comments")
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("User is not authenticated")
		return &postService.CommentsResponse{}, err
	}
	err = s.roleHavePermission(role, "post_getAll")
	if err != nil {
		Log.Warn("User doesn't have permission to get requests")
		return &postService.CommentsResponse{}, err
	}

	return s.postClient.GetAllCommentsRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllCommentsFromPostRequest(ctx context.Context, in *postService.PostCommentsRequest) (*postService.CommentsResponse, error) {
	Log.Info("User with id: " + in.LoggedUserId + " getting all comment for post with id: " + in.PostId)
	in.LoggedUserId = getUserIdFromJwt(ctx)
	return s.postClient.GetAllCommentsFromPostRequest(ctx, in)
}

func (s *PostGatewayStruct) CreateCommentRequest(ctx context.Context, in *postService.CommentRequest) (*postService.CommentResponse, error) {
	Log.Info("User with id: " + in.LoggedUserId + " getting comment for post with id: " + in.Id)
	err := checkValue(in.String())
	if err != nil {
		Log.Warn("Input possibly contains malicious data")
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("User is not authenticated")
		return &postService.CommentResponse{}, err
	}
	err = s.roleHavePermission(role, "post_write")
	if err != nil {
		Log.Warn("User doesn't have permission to get requests")
		return &postService.CommentResponse{}, err
	}
	in.LoggedUserId = getUserIdFromJwt(ctx)

	return s.postClient.CreateCommentRequest(ctx, in)
}

func (s *PostGatewayStruct) DeleteCommentRequest(ctx context.Context, in *postService.CommentIdRequest) (*postService.EmptyRequest, error) {
	Log.Info("User with id: " + in.LoggedUserId + " deleting comment with id: " + in.Id)
	err := checkValue(in.String())
	if err != nil {
		Log.Warn("Input possibly contains malicious data")
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("User is not authenticated")
		return &postService.EmptyRequest{}, err
	}
	err = s.roleHavePermission(role, "post_delete")
	if err != nil {
		Log.Warn("User doesn't have permission to get requests")
		return &postService.EmptyRequest{}, err
	}
	in.LoggedUserId = getUserIdFromJwt(ctx)
	return s.postClient.DeleteCommentRequest(ctx, in)
}

func (s *PostGatewayStruct) GetReactionRequest(ctx context.Context, in *postService.ReactionIdRequest) (*postService.ReactionResponse, error) {
	Log.Info("User with id: " + in.LoggedUserId + " getting reaction with id: " + in.Id)
	in.LoggedUserId = getUserIdFromJwt(ctx)
	return s.postClient.GetReactionRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllReactionsRequest(ctx context.Context, in *postService.EmptyRequest) (*postService.ReactionsResponse, error) {
	Log.Info("Getting all reactions")
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("User is not authenticated")
		return nil, err
	}
	err = s.roleHavePermission(role, "post_getAll")
	if err != nil {
		Log.Warn("User doesn't have permission to get requests")
		return nil, err
	}

	return s.postClient.GetAllReactionsRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllReactionsFromPostRequest(ctx context.Context, in *postService.PostReactionRequest) (*postService.ReactionsResponse, error) {
	Log.Info("User with id: " + in.LoggedUserId + " getting reaction from post with id: " + in.PostId)
	return s.postClient.GetAllReactionsFromPostRequest(ctx, in)
}

func (s *PostGatewayStruct) CreateReactionRequest(ctx context.Context, in *postService.ReactionRequest) (*postService.ReactionResponse, error) {
	Log.Info("User with id: " + in.LoggedUserId + " creating reaction on post with id: " + in.PostId)
	err := checkValue(in.String())
	if err != nil {
		Log.Warn("Input possibly contains malicious data")
		return nil, err
	}
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("User is not authenticated")
		return nil, err
	}
	err = s.roleHavePermission(role, "post_write")
	if err != nil {
		Log.Warn("User doesn't have permission to get requests")
		return nil, err
	}
	in.LoggedUserId = getUserIdFromJwt(ctx)

	return s.postClient.CreateReactionRequest(ctx, in)
}

func (s *PostGatewayStruct) DeleteReactionRequest(ctx context.Context, in *postService.ReactionIdRequest) (*postService.EmptyRequest, error) {
	Log.Info("User with id: " + in.LoggedUserId + " deleting reaction with id: " + in.Id)
	role, err := s.isUserAuthenticated(ctx)
	if err != nil {
		Log.Warn("User is not authenticated")
		return nil, err
	}
	err = s.roleHavePermission(role, "post_delete")
	if err != nil {
		Log.Warn("User doesn't have permission to get requests")
		return nil, err
	}
	in.LoggedUserId = getUserIdFromJwt(ctx)

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
