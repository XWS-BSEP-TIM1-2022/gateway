package api

import (
	"errors"
	"fmt"
	"gateway/startup/config"
	jobService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/job"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

type JobGatewayStruct struct {
	jobService.UnimplementedJobServiceServer
	config     *config.Config
	jobClient  jobService.JobServiceClient
	userClient userService.UserServiceClient
}

func NewJobGateway(c *config.Config) *JobGatewayStruct {
	return &JobGatewayStruct{
		config:     c,
		jobClient:  services.NewJobClient(fmt.Sprintf("%s:%s", c.JobServiceHost, c.JobServicePort)),
		userClient: services.NewUserClient(fmt.Sprintf("%s:%s", c.UserServiceHost, c.UserServicePort)),
	}
}

func (s *JobGatewayStruct) GetRequest(ctx context.Context, in *jobService.JobIdRequest) (*jobService.GetResponse, error) {
	return s.jobClient.GetRequest(ctx, in)
}

func (s *JobGatewayStruct) GetAllRequest(ctx context.Context, in *jobService.EmptyRequest) (*jobService.JobsResponse, error) {
	return s.jobClient.GetAllRequest(ctx, in)
}

func (s *JobGatewayStruct) PostRequest(ctx context.Context, in *jobService.UserRequest) (*jobService.GetResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	jwt := md.Get("Authorization")
	if jwt == nil {
		return nil, errors.New("unauthorized")
	}
	userId, err := s.userClient.IsApiTokenValid(ctx, &userService.AuthRequest{Token: jwt[0]})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	in.Job.UserId = userId.UserId
	return s.jobClient.PostRequest(ctx, in)
}

func (s *JobGatewayStruct) DeleteRequest(ctx context.Context, in *jobService.JobIdRequest) (*jobService.EmptyRequest, error) {
	return s.jobClient.DeleteRequest(ctx, in)
}

func (s *JobGatewayStruct) SearchJobsRequest(ctx context.Context, in *jobService.SearchRequest) (*jobService.JobsResponse, error) {
	return s.jobClient.SearchJobsRequest(ctx, in)
}
