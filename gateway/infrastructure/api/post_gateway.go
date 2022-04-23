package api

import (
	"context"
	"fmt"
	"gateway/infrastructure/services"
	"gateway/startup/config"
	postService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/post"
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
