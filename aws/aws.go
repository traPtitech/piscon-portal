package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	ImageId      = "ami-e7527ed7"            //TODO
	InstanceType = types.InstanceTypeT2Micro //TODO
)

type AwsClient struct {
	c *ec2.Client
}

func New() *AwsClient {
	a := &AwsClient{}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	client := ec2.NewFromConfig(cfg)
	a.c = client
	return a
}

func (a *AwsClient) CreateInstances(c context.Context, name string, num int32, subnetId string, privateIp string) error {
	i := &ec2.RunInstancesInput{
		ImageId:          aws.String(ImageId),
		InstanceType:     InstanceType,
		MinCount:         &num,
		MaxCount:         &num,
		SubnetId:         &subnetId,
		PrivateIpAddress: &privateIp,
	}
	res, err := a.c.RunInstances(c, i)
	if err != nil {
		return err
	}

	err = a.CreateTags(c, *res.Instances[0].InstanceId, "Name", name)
	if err != nil {
		return err
	}
	return nil
}

func (a *AwsClient) CreateTags(c context.Context, instanceId string, key string, value string) error {
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

func (a *AwsClient) DeleteInstances(c context.Context, instanceId string) error {
	i := &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceId},
	}
	_, err := a.c.TerminateInstances(c, i)
	if err != nil {
		return err
	}
	return nil
}

func (a *AwsClient) StartInstances(c context.Context, instanceId string) error {
	i := &ec2.StartInstancesInput{
		InstanceIds: []string{instanceId},
	}
	_, err := a.c.StartInstances(c, i)
	if err != nil {
		return err
	}
	return nil
}

func (a *AwsClient) StopInstances(c context.Context, instanceId string) error {
	i := &ec2.StopInstancesInput{
		InstanceIds: []string{instanceId},
	}
	_, err := a.c.StopInstances(c, i)
	if err != nil {
		return err
	}
	return nil
}

func (a *AwsClient) GetInstancesInfo(c context.Context, instanceId string) error {
	return a.c.DescribeInstances(c, input)
}
