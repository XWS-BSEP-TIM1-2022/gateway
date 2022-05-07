package api

import (
	"context"
	"fmt"
	"gateway/startup/config"
	postService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/post"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
)

type PostGatewayStruct struct {
	postService.UnimplementedPostServiceServer
	config     *config.Config
	postClient postService.PostServiceClient
}

func NewPostGateway(c *config.Config) *PostGatewayStruct {
	return &PostGatewayStruct{
		config:     c,
		postClient: services.NewPostClient(fmt.Sprintf("%s:%s", c.PostServiceHost, c.PostServicePort)),
	}
}

func (s *PostGatewayStruct) GetRequest(ctx context.Context, in *postService.PostIdRequest) (*postService.PostResponse, error) {
	return s.postClient.GetRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllRequest(ctx context.Context, in *postService.EmptyRequest) (*postService.PostsResponse, error) {
	return s.postClient.GetAllRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllFromUserRequest(ctx context.Context, in *postService.UserPostsRequest) (*postService.PostsResponse, error) {
	return s.postClient.GetAllFromUserRequest(ctx, in)
}

func (s *PostGatewayStruct) CreateRequest(ctx context.Context, in *postService.PostRequest) (*postService.PostResponse, error) {
	return s.postClient.CreateRequest(ctx, in)
}

func (s *PostGatewayStruct) DeleteRequest(ctx context.Context, in *postService.PostIdRequest) (*postService.EmptyRequest, error) {
	return s.postClient.DeleteRequest(ctx, in)
}

func (s *PostGatewayStruct) GetCommentRequest(ctx context.Context, in *postService.CommentIdRequest) (*postService.CommentResponse, error) {
	return s.postClient.GetCommentRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllCommentsRequest(ctx context.Context, in *postService.EmptyRequest) (*postService.CommentsResponse, error) {
	return s.postClient.GetAllCommentsRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllCommentsFromPostRequest(ctx context.Context, in *postService.PostCommentsRequest) (*postService.CommentsResponse, error) {
	return s.postClient.GetAllCommentsFromPostRequest(ctx, in)
}

func (s *PostGatewayStruct) CreateCommentRequest(ctx context.Context, in *postService.CommentRequest) (*postService.CommentResponse, error) {
	return s.postClient.CreateCommentRequest(ctx, in)
}

func (s *PostGatewayStruct) DeleteCommentRequest(ctx context.Context, in *postService.CommentIdRequest) (*postService.EmptyRequest, error) {
	return s.postClient.DeleteCommentRequest(ctx, in)
}

func (s *PostGatewayStruct) GetReactionRequest(ctx context.Context, in *postService.ReactionIdRequest) (*postService.ReactionResponse, error) {
	return s.postClient.GetReactionRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllReactionsRequest(ctx context.Context, in *postService.EmptyRequest) (*postService.ReactionsResponse, error) {
	return s.postClient.GetAllReactionsRequest(ctx, in)
}

func (s *PostGatewayStruct) GetAllReactionsFromPostRequest(ctx context.Context, in *postService.PostReactionRequest) (*postService.ReactionsResponse, error) {
	return s.postClient.GetAllReactionsFromPostRequest(ctx, in)
}

func (s *PostGatewayStruct) CreateReactionRequest(ctx context.Context, in *postService.ReactionRequest) (*postService.ReactionResponse, error) {
	return s.postClient.CreateReactionRequest(ctx, in)
}

func (s *PostGatewayStruct) DeleteReactionRequest(ctx context.Context, in *postService.ReactionIdRequest) (*postService.EmptyRequest, error) {
	return s.postClient.DeleteReactionRequest(ctx, in)
}
