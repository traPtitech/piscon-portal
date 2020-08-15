package conoha

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/startstop"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/pagination"
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

func (c *ConohaClient) MakeInstance(name, privateIP string) error {
	eo := gophercloud.EndpointOpts{
		Type:   "compute",
		Region: "tyo2",
	}
	compute, err := openstack.NewComputeV2(c.client, eo)
	if err != nil {
		panic(err)
	}

	startUpScript := fmt.Sprintf(`#!/bin/sh

sed -e "s/PermitRootLogin yes/PermitRootLogin no/g" -i /etc/ssh/sshd_config
sed -e "s/#PermitRootLogin yes/PermitRootLogin no/g" -i /etc/ssh/sshd_config
sed -e "s/#PermitRootLogin no/PermitRootLogin no/g" -i /etc/ssh/sshd_config

cat <<EOF>/etc/netplan/11-privatenetwork.yaml
network:
    ethernets:
        eth1:
            addresses: [%s/21]
            dhcp4: false
    version: 2
EOF

systemctl restart sshd	
	`, privateIP)

	copts := servers.CreateOpts{
		Name:      name,
		ImageRef:  os.Getenv("CONOHA_IMAGE_REF"),
		FlavorRef: os.Getenv("CONOHA_IMAGE_FLAVOR"),
		Metadata: map[string]string{
			"instance_name_tag": name,
		},
		SecurityGroups: []string{"default", "gncs-ipv4-all", "gncs-ipv6-all"},
		UserData:       []byte(startUpScript),
	}
	log.Println(startUpScript)
	r := servers.Create(compute, copts)
	if r.Err != nil {
		return r.Err
	}
	return nil
}

func (c *ConohaClient) DeleteInstance(name string) error {
	eo := gophercloud.EndpointOpts{
		Type:   "compute",
		Region: "tyo2",
	}
	compute, err := openstack.NewComputeV2(c.client, eo)
	if err != nil {
		panic(err)
	}
	instance, err := c.GetInstanceInfo(name)
	// ports := server.
	err = servers.Delete(compute, instance.ID).ExtractErr()
	if err != nil {
		return err
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
		Region: "tyo2",
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
		Region: "tyo2",
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

func (c *ConohaClient) ShutdownInstance(instanceName string) error {
	instance, err := c.GetInstanceInfo(instanceName)
	if err != nil {
		return err
	}

	eo := gophercloud.EndpointOpts{
		Type:   "compute",
		Region: "tyo2",
	}
	compute, err := openstack.NewComputeV2(c.client, eo)
	if err != nil {
		panic(err)
	}

	err = startstop.Stop(compute, instance.ID).ExtractErr()
	if err != nil {
		return err
	}
	return nil
}

func (c *ConohaClient) StartInstance(instanceName string) error {
	instance, err := c.GetInstanceInfo(instanceName)
	if err != nil {
		return err
	}

	eo := gophercloud.EndpointOpts{
		Type:   "compute",
		Region: "tyo2",
	}
	compute, err := openstack.NewComputeV2(c.client, eo)
	if err != nil {
		panic(err)
	}

	err = startstop.Start(compute, instance.ID).ExtractErr()
	if err != nil {
		return err
	}
	return nil
}

func (c *ConohaClient) AttachPrivateNetwork(instanceName, networkID, privateIP string) error {
	instance, err := c.GetInstanceInfo(instanceName)
	if err != nil {
		return err
	}
	port := c.CreatePorts(privateIP, networkID)

	eo := gophercloud.EndpointOpts{
		Type:   "compute",
		Region: "tyo2",
	}
	compute, err := openstack.NewComputeV2(c.client, eo)
	if err != nil {
		panic(err)
	}
	attachOpts := attachinterfaces.CreateOpts{
		// NetworkID: networkID,
		PortID: port.ID,
	}
	log.Println("create attachinterface")
	_, err = attachinterfaces.Create(compute, instance.ID, attachOpts).Extract()

	if err != nil {
		panic(err)
	}
	return nil
}

func (c *ConohaClient) CreatePorts(privateIP, networkID string) *ports.Port {
	log.Println("Create ports")
	eo := gophercloud.EndpointOpts{
		Type:   "network",
		Region: "tyo2",
	}
	networkClient, err := openstack.NewNetworkV2(c.client, eo)
	if err != nil {
		panic(err)
	}

	createOpts := ports.CreateOpts{
		Name: "private-port",
		// AdminStateUp: &asu,
		NetworkID: os.Getenv("CONOHA_NETWORK_ID"),
		FixedIPs: []ports.IP{
			{SubnetID: os.Getenv("CONOHA_SUBNET_ID"), IPAddress: privateIP},
		},
	}

	port, err := ports.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}
	return port
}
