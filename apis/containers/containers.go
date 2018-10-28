package containers

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/containerinstance/mgmt/containerinstance"
	"github.com/iphilpot/flare/apis/config"
	"github.com/iphilpot/flare/apis/errors"
	"github.com/iphilpot/flare/apis/iam"
	"github.com/iphilpot/flare/apis/logger"
)

func getContainerInstanceClient() containerinstance.ContainerGroupsClient {
	authorizer := iam.GetAuthorizerFromEnvironment()
	c := config.GetConfig()
	containerClient := containerinstance.NewContainerGroupsClient(c.AzureSubscriptionID)
	containerClient.Authorizer = authorizer
	return containerClient
}

func defineContainerProperties(imageName string, port int32, memory, cpu float64) containerinstance.ContainerProperties {
	return containerinstance.ContainerProperties{
		Image: &imageName,
		Ports: &[]containerinstance.ContainerPort{
			containerinstance.ContainerPort{
				Protocol: containerinstance.ContainerNetworkProtocolTCP,
				Port:     &port,
			},
		},
		Resources: &containerinstance.ResourceRequirements{
			Requests: &containerinstance.ResourceRequests{
				MemoryInGB: &memory,
				CPU:        &cpu,
			},
		},
	}
}

func defineContainerGroup(location, dnsName string, port int32, containers []Container) containerinstance.ContainerGroup {
	var containerSet []containerinstance.Container
	for i := range containers {
		currContainer := containerinstance.Container{
			Name:                &containers[i].Name,
			ContainerProperties: &containers[i].Properties,
		}
		containerSet = append(containerSet, currContainer)
	}
	return containerinstance.ContainerGroup{
		Location: &location,
		ContainerGroupProperties: &containerinstance.ContainerGroupProperties{
			Containers:    &containerSet,
			RestartPolicy: containerinstance.Always,
			OsType:        containerinstance.Linux,
			IPAddress: &containerinstance.IPAddress{
				Ports: &[]containerinstance.Port{
					containerinstance.Port{
						Protocol: containerinstance.TCP,
						Port:     &port,
					},
				},
				Type:         containerinstance.Public,
				DNSNameLabel: &dnsName,
			},
		},
	}
}

// Container - Desctibes the containers you want to create
type Container struct {
	Name       string
	ImageName  string
	Port       int32
	CPU        float64
	Memory     float64
	Properties containerinstance.ContainerProperties
}

// ContainerGroup - Describes the group of containers to be created
type ContainerGroup struct {
	GroupName         string
	ResourceGroupName string
	Location          string
	DNSName           string
	Port              int32
	Containers        []Container
}

// CreateContainerGroup - Creates a container group
func (cg ContainerGroup) CreateContainerGroup(ctx context.Context) {
	client := getContainerInstanceClient()
	for i := range cg.Containers {
		cg.Containers[i].Properties = defineContainerProperties(cg.Containers[i].ImageName,
			cg.Containers[i].Port, cg.Containers[i].Memory, cg.Containers[i].CPU)
	}
	containerGroup := defineContainerGroup(cg.Location, cg.DNSName, cg.Port, cg.Containers)
	future, err := client.CreateOrUpdate(ctx, cg.ResourceGroupName, cg.GroupName, containerGroup)
	errors.HandleError(err)
	err = future.WaitForCompletion(ctx, client.Client)
	errors.HandleError(err)
	result, _ := future.Result(client)
	logger.PrintAndLog(fmt.Sprintf("Container Group: %s has been created.", *result.Name))
}
