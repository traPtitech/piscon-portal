package conoha

import (
	"errors"
	"fmt"
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/images"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
)

type ConohaClient struct {
	client *gophercloud.ProviderClient
}

func New(opts gophercloud.AuthOptions) *ConohaClient {
	c := &ConohaClient{}
	client, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		panic(err)
	}
	c.client = client
	return c
}

func (c *ConohaClient) MakeInstance(name, pass string) error {
	eo := gophercloud.EndpointOpts{
		Type:   "compute",
		Region: "tyo1",
	}
	compute, err := openstack.NewComputeV2(c.client, eo)
	if err != nil {
		panic(err)
	}

	startUpScript := fmt.Sprintf(`#!/bin/sh

useradd -m -G sudo -s /bin/bash isucon
echo isucon:%s | /usr/sbin/chpasswd

sed -e "s/PermitRootLogin yes/PermitRootLogin no/g" -i /etc/ssh/sshd_config
sed -e "s/#PermitRootLogin yes/PermitRootLogin no/g" -i /etc/ssh/sshd_config
sed -e "s/#PermitRootLogin no/PermitRootLogin no/g" -i /etc/ssh/sshd_config

systemctl restart sshd	
	`, pass)

	copts := servers.CreateOpts{
		Name:      "isucon",
		ImageRef:  os.Getenv("CONOHA_IMAGE_REF"),
		FlavorRef: os.Getenv("CONOHA_IMAGE_FLAVOR"),
		Metadata: map[string]string{
			"instance_name_tag": name,
		},
		SecurityGroups: []string{os.Getenv("CONOHA_SECURITY_GROUP")},
		UserData:       []byte(startUpScript),
	}
	r := servers.Create(compute, copts)
	if r.Err != nil {
		return r.Err
	}
	return nil
}

func (c *ConohaClient) GetInstanceInfo(instanceName string) (*servers.Server, error) {
	list, err := c.InstanceList()
	if err != nil {
		return nil, err
	}
	for _, server := range list {
		fmt.Println(server.Metadata["instance_name_tag"])
		if server.Metadata["instance_name_tag"] == instanceName {
			return &server, nil
		}
	}
	return nil, errors.New("Not Found")
}

func (c *ConohaClient) InstanceList() ([]servers.Server, error) {
	eo := gophercloud.EndpointOpts{
		Type:   "compute",
		Region: "tyo1",
	}
	compute, err := openstack.NewComputeV2(c.client, eo)
	if err != nil {
		panic(err)
	}

	opts := servers.ListOpts{}

	pager := servers.List(compute, opts)

	list := make([]servers.Server, 0)

	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, _ := servers.ExtractServers(page)
		list = append(list, serverList...)
		return true, nil
	})
	return list, err
}

func (c *ConohaClient) ImageList() ([]images.Image, error) {
	eo := gophercloud.EndpointOpts{
		Type:   "compute",
		Region: "tyo1",
	}
	compute, err := openstack.NewComputeV2(c.client, eo)
	if err != nil {
		panic(err)
	}

	opts := images.ListOpts{}

	pager := images.ListDetail(compute, opts)

	list := make([]images.Image, 0)

	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		imageList, _ := images.ExtractImages(page)
		list = append(list, imageList...)
		return true, nil
	})
	return list, err
}
