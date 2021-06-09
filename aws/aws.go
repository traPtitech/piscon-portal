package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2CreateInstanceAPI interface {
	RunInstances(ctx context.Context, params *ec2.RunInstancesInput, optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)
	CreateTags(ctx context.Context, params *ec2.CreateTagsInput, optFns ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)
}

type EC2DeleteInstanceAPI interface {
	TerminateInstances(ctx context.Context, params *ec2.TerminateInstancesInput, optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}

type EC2InstanceAPI interface {
	StartInstances(ctx context.Context, params *ec2.StartInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error)
	StopInstances(ctx context.Context, params *ec2.StopInstancesInput, optFns ...func(*ec2.Options)) (*ec2.StopInstancesOutput, error)
}

type AwsClient struct {
	Client *ec2.Client
}

func New() *AwsClient {
	a := &AwsClient{}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	client := ec2.NewFromConfig(cfg)
	a.Client = client
	return a
}

func (a *AwsClient) CreateInstance(c context.Context, api EC2CreateInstanceAPI, input *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {
	return api.RunInstances(c, input)
}

func (a *AwsClient) CreateTags(c context.Context, api EC2CreateInstanceAPI, input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return api.CreateTags(c, input)
}

func (a *AwsClient) DeleteInstance(c context.Context, api EC2DeleteInstanceAPI, input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	return api.TerminateInstances(c, input)
}
func (a *AwsClient) StartInstances(c context.Context, api EC2InstanceAPI, input *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
	return api.StartInstances(c, input)
}
func (a *AwsClient) StopInstances(c context.Context, api EC2InstanceAPI, input *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
	return api.StopInstances(c, input)
}
