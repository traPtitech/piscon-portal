package aws

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/traPtitech/piscon-portal/model"
)

const (
	ImageId      = "ami-e7527ed7"            //TODO
	InstanceType = types.InstanceTypeT2Micro //TODO
)

var (
	defaultInstanceNum = int32(1)
	InstanceNameKey    = "Name"
	statusmap          = map[string]string{
		string(types.InstanceStateNamePending):    model.STARTING, //TODO buildingと被っている
		string(types.InstanceStateNameRunning):    model.ACTIVE,
		string(types.InstanceStateNameTerminated): model.NOT_EXIST,
		string(types.InstanceStateNameStopping):   model.SHUTDOWNING,
		string(types.InstanceStateNameStopped):    model.SHUTOFF,
	}
)

type AwsClient struct {
	c *ec2.Client
}

func New(cfg aws.Config) (*AwsClient, error) {
	a := &AwsClient{}
	client := ec2.NewFromConfig(cfg)
	a.c = client
	return a, nil
}

func CreateDefaultConfig() (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     os.Getenv("ACCESS_ID"),
					SecretAccessKey: os.Getenv("ACCESS_SECRET_KEY"),
				},
			},
		),
	)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (a *AwsClient) CreateInstance(c context.Context, name string, privateIp string) (*string, error) {
	subnetId := os.Getenv("AWS_SUBNET_ID")
	i := &ec2.RunInstancesInput{
		ImageId:          aws.String(ImageId),
		InstanceType:     InstanceType,
		MinCount:         &defaultInstanceNum,
		MaxCount:         &defaultInstanceNum,
		SubnetId:         &subnetId,
		PrivateIpAddress: &privateIp,
	}
	res, err := a.c.RunInstances(c, i)
	if err != nil {
		return nil, err
	}

	err = a.CreateTag(c, *res.Instances[0].InstanceId, "Name", name)
	if err != nil {
		return nil, err
	}
	return res.Instances[0].InstanceId, nil
}

func (a *AwsClient) CreateTag(c context.Context, instanceId string, key string, value string) error {
	i := &ec2.CreateTagsInput{
		Resources: []string{instanceId},
		Tags: []types.Tag{
			{
				Key:   &key,
				Value: &value,
			},
		},
	}
	_, err := a.c.CreateTags(c, i)
	if err != nil {
		return err
	}
	return nil
}

func (a *AwsClient) DeleteInstance(c context.Context, instanceId string) error {
	i := &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceId},
	}
	_, err := a.c.TerminateInstances(c, i)
	if err != nil {
		return err
	}
	return nil
}

func (a *AwsClient) StartInstance(c context.Context, instanceId string) error {
	i := &ec2.StartInstancesInput{
		InstanceIds: []string{instanceId},
	}
	_, err := a.c.StartInstances(c, i)
	if err != nil {
		return err
	}
	return nil
}

func (a *AwsClient) StopInstance(c context.Context, instanceId string) error {
	i := &ec2.StopInstancesInput{
		InstanceIds: []string{instanceId},
	}
	_, err := a.c.StopInstances(c, i)
	if err != nil {
		return err
	}
	return nil
}

func (a *AwsClient) GetInstanceInfo(c context.Context, instanceName string) (*model.Instance, error) {
	i := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   &InstanceNameKey,
				Values: []string{instanceName},
			},
		},
	}
	res, err := a.c.DescribeInstances(c, i)
	if err != nil {
		return nil, err
	}
	instance := &model.Instance{
		GlobalIPAddress:  *res.Reservations[0].Instances[0].PublicIpAddress,
		PrivateIPAddress: *res.Reservations[0].Instances[0].PrivateIpAddress,
		Status:           statusmap[string(res.Reservations[0].Instances[0].State.Name)],
	}
	return instance, nil
}
