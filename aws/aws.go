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
	InstanceType = types.InstanceTypeT2Micro //TODO                  //TODO
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
	var k string
	k = "Name"
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
	ti := &ec2.CreateTagsInput{
		Resources: []string{*res.Instances[0].InstanceId},
		Tags: []types.Tag{
			{
				Key:   &k,
				Value: &name,
			},
		},
	}
	_, err = a.CreateTags(c, ti)
	if err != nil {
		return err
	}
	return nil
}

func (a *AwsClient) CreateTags(c context.Context, input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return a.c.CreateTags(c, input)
}

func (a *AwsClient) DeleteInstances(c context.Context, input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	return a.c.TerminateInstances(c, input)
}

func (a *AwsClient) StartInstances(c context.Context, input *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
	return a.c.StartInstances(c, input)
}

func (a *AwsClient) StopInstances(c context.Context, input *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
	return a.c.StopInstances(c, input)
}

func (a *AwsClient) GetInstancesInfo(c context.Context, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return a.c.DescribeInstances(c, input)
}
