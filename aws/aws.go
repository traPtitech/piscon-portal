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
	imageId      = string("ami-03b1b78bb1da5122f") // isucon競技用サーバーのAMI
	InstanceType = types.InstanceTypeT2Medium      // isuconサーバーの種類(競技ごとにスペックが違う)
	region       = string("ap-northeast-1")        // isuconサーバーのリージョン
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
					AccessKeyID:     os.Getenv("AWS_ACCESS_KEY"),
					SecretAccessKey: os.Getenv("AWS_ACCESS_SECRET"),
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
		}, {
			Key:   aws.String("type"),
			Value: aws.String("PISCON"),
		}},
	}
	startUpScript := fmt.Sprintf(`#!/bin/sh
useradd -m isucon
echo "%s\n%s\n" | passwd isucon
usermod -G sudo isucon
sed -e "s/PasswordAuthentication no/PasswordAuthentication yes/g" -i /etc/ssh/sshd_config
systemctl restart sshd
echo "server {" > /etc/nginx/sites-available/isucari.conf
echo "	# listen 443 ssl;" >> /etc/nginx/sites-available/isucari.conf
echo "	# server_name isucon9.catatsuy.org;" >> /etc/nginx/sites-available/isucari.conf
echo "" >> /etc/nginx/sites-available/isucari.conf
echo "	# ssl_certificate //etc/nginx/sites-available/isucari.confssl/fullchain.pem;" >> /etc/nginx/sites-available/isucari.conf
echo "	# ssl_certificate_key //etc/nginx/sites-available/isucari.confssl/privkey.pem;" >> /etc/nginx/sites-available/isucari.conf
echo "" >> /etc/nginx/sites-available/isucari.conf
echo "	location / {" >> /etc/nginx/sites-available/isucari.conf
echo "			proxy_set_header Host $http_host;" >> /etc/nginx/sites-available/isucari.conf
echo "			proxy_pass http://127.0.0.1:8000;" >> /etc/nginx/sites-available/isucari.conf
echo "	}" >> /etc/nginx/sites-available/isucari.conf
echo "}" >> /etc/nginx/sites-available/isucari.conf
	`, pwd, pwd)
	enc := base64.StdEncoding.EncodeToString([]byte(startUpScript))
	nispec := types.InstanceNetworkInterfaceSpecification{
		AssociatePublicIpAddress: aws.Bool(true),
		DeleteOnTermination:      aws.Bool(true),
		DeviceIndex:              aws.Int32(0),
		SubnetId:                 &subnetId,
		PrivateIpAddress:         &privateIp,
		Groups:                   []string{os.Getenv("AWS_SECURITY_GROUP_ID")},
	}
	i := &ec2.RunInstancesInput{
		ImageId:           aws.String(imageId),
		InstanceType:      InstanceType,
		MinCount:          &defaultInstanceNum,
		MaxCount:          &defaultInstanceNum,
		TagSpecifications: []types.TagSpecification{tspec},
		NetworkInterfaces: []types.InstanceNetworkInterfaceSpecification{nispec},
		KeyName:           aws.String("piscon-portal"),
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
		GlobalIPAddress:  aws.ToString(res.Reservations[0].Instances[0].PublicIpAddress),
		PrivateIPAddress: aws.ToString(res.Reservations[0].Instances[0].PrivateIpAddress),
		Status:           statusmap[string(res.Reservations[0].Instances[0].State.Name)],
	}
	fmt.Println(res.Reservations[0].Instances[0].State.Name)
	fmt.Println(statusmap[string(res.Reservations[0].Instances[0].State.Name)])
	return instance, nil
}
