package api

import (
	"fmt"
	"gateway/startup/config"
	jobService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/job"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
	"golang.org/x/net/context"
)

type JobGatewayStruct struct {
	jobService.UnimplementedJobServiceServer
	config    *config.Config
	jobClient jobService.JobServiceClient
}

func NewJobGateway(c *config.Config) *JobGatewayStruct {
	return &JobGatewayStruct{
		config:    c,
		jobClient: services.NewJobClient(fmt.Sprintf("%s:%s", c.JobServiceHost, c.JobServicePort)),
	}
}

func (s *JobGatewayStruct) GetRequest(ctx context.Context, in *jobService.JobIdRequest) (*jobService.GetResponse, error) {
	return s.jobClient.GetRequest(ctx, in)
}

func (s *JobGatewayStruct) GetAllRequest(ctx context.Context, in *jobService.EmptyRequest) (*jobService.JobsResponse, error) {
	return s.jobClient.GetAllRequest(ctx, in)
}

func (s *JobGatewayStruct) PostRequest(ctx context.Context, in *jobService.UserRequest) (*jobService.GetResponse, error) {
	return s.jobClient.PostRequest(ctx, in)
}

func (s *JobGatewayStruct) DeleteRequest(ctx context.Context, in *jobService.JobIdRequest) (*jobService.EmptyRequest, error) {
	return s.jobClient.DeleteRequest(ctx, in)
}

func (s *JobGatewayStruct) SearchJobsRequest(ctx context.Context, in *jobService.SearchRequest) (*jobService.JobsResponse, error) {
	return s.jobClient.SearchJobsRequest(ctx, in)
}
