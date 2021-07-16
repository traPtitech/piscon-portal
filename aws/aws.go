package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/traPtitech/piscon-portal/model"
)

const (
	imageId      = string("ami-03b1b78bb1da5122f") //TODO
	InstanceType = types.InstanceTypeT2Micro       //TODO
	region       = string("ap-northeast-1")
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

type Config aws.Config //TODO 苦肉の策、いい感じに分離したい

type AwsClient struct {
	c *ec2.Client
}

func New(cfg Config) (*AwsClient, error) {
	a := &AwsClient{}
	client := ec2.NewFromConfig(aws.Config(cfg))
	a.c = client
	return a, nil
}

func CreateDefaultConfig() (*Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     os.Getenv("ACCESS_ID"),
					SecretAccessKey: os.Getenv("ACCESS_SECRET_KEY"),
				},
			},
		), config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	res := Config(cfg)
	return &res, nil
}

func (a *AwsClient) CreateInstance(name string, privateIp string, pwd string) (*string, error) {
	subnetId := os.Getenv("AWS_SUBNET_ID")
	tspec := types.TagSpecification{
		ResourceType: types.ResourceTypeInstance,
		Tags: []types.Tag{{
			Key:   aws.String("Name"),
			Value: &name,
		}},
	}
	startUpScript := fmt.Sprintf(`#!/bin/sh
useradd -m piscon
echo "%s\n%s\n" | passwd piscon
usermod -G sudo piscon
sed -e "s/PasswordAuthentication no/PasswordAuthentication yes/g" -i /etc/ssh/sshd_config
systemctl restart sshd	
	`, pwd, pwd)
	fmt.Printf(startUpScript)
	enc := base64.StdEncoding.EncodeToString([]byte(startUpScript))
	nispec := types.InstanceNetworkInterfaceSpecification{
		AssociatePublicIpAddress: aws.Bool(true),
		DeleteOnTermination:      aws.Bool(true),
		DeviceIndex:              aws.Int32(0),
		SubnetId:                 &subnetId,
		PrivateIpAddress:         &privateIp,
	}
	i := &ec2.RunInstancesInput{
		ImageId:           aws.String(imageId),
		InstanceType:      InstanceType,
		MinCount:          &defaultInstanceNum,
		MaxCount:          &defaultInstanceNum,
		TagSpecifications: []types.TagSpecification{tspec},
		NetworkInterfaces: []types.InstanceNetworkInterfaceSpecification{nispec},
		KeyName:           aws.String("isucon"),
		UserData:          aws.String(enc),
	}
	res, err := a.c.RunInstances(context.TODO(), i)
	if err != nil {
		return nil, err
	}
	return res.Instances[0].InstanceId, nil
}

func (a *AwsClient) CreateTag(instanceId string, key string, value string) error {
	i := &ec2.CreateTagsInput{
		Resources: []string{instanceId},
		Tags: []types.Tag{
			{
				Key:   &key,
				Value: &value,
			},
		},
	}
	_, err := a.c.CreateTags(context.TODO(), i)
	if err != nil {
		return err
	}
	return nil
}

func (a *AwsClient) DeleteInstance(instanceId string) error {
	i := &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceId},
	}
	_, err := a.c.TerminateInstances(context.TODO(), i)
	if err != nil {
		return err
	}
	return nil
}

func (a *AwsClient) StartInstance(instanceId string) error {
	i := &ec2.StartInstancesInput{
		InstanceIds: []string{instanceId},
	}
	_, err := a.c.StartInstances(context.TODO(), i)
	if err != nil {
		return err
	}
	return nil
}

func (a *AwsClient) StopInstance(instanceId string) error {
	i := &ec2.StopInstancesInput{
		InstanceIds: []string{instanceId},
	}
	_, err := a.c.StopInstances(context.TODO(), i)
	if err != nil {
		return err
	}
	return nil
}

func (a *AwsClient) GetInstanceInfo(id string) (*model.Instance, error) {
	i := &ec2.DescribeInstancesInput{
		InstanceIds: []string{id},
	}
	res, err := a.c.DescribeInstances(context.TODO(), i)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	instance := &model.Instance{
		GlobalIPAddress:  *res.Reservations[0].Instances[0].PublicIpAddress,
		PrivateIPAddress: *res.Reservations[0].Instances[0].PrivateIpAddress,
		Status:           statusmap[string(res.Reservations[0].Instances[0].State.Name)],
	}
	fmt.Println(res.Reservations[0].Instances[0].State.Name)
	fmt.Println(statusmap[string(res.Reservations[0].Instances[0].State.Name)])
	return instance, nil
}
