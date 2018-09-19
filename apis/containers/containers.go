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

var (
	ipType = "Public"
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

func defineContainerGroup(containerName, location, dnsName string, port int32, containerProperties containerinstance.ContainerProperties) containerinstance.ContainerGroup {
	return containerinstance.ContainerGroup{
		Location: &location,
		ContainerGroupProperties: &containerinstance.ContainerGroupProperties{
			Containers: &[]containerinstance.Container{
				containerinstance.Container{
					Name:                &containerName,
					ContainerProperties: &containerProperties,
				},
			},
			RestartPolicy: containerinstance.Always,
			OsType:        containerinstance.Linux,
			IPAddress: &containerinstance.IPAddress{
				Ports: &[]containerinstance.Port{
					containerinstance.Port{
						Protocol: containerinstance.TCP,
						Port:     &port,
					},
				},
				Type:         &ipType,
				DNSNameLabel: &dnsName,
			},
		},
	}
}

// CreateContainer - Creates a container
func CreateContainer(ctx context.Context, resourceGroupName, containerGroupName, location, dnsName string) {
	client := getContainerInstanceClient()
	cProps := defineContainerProperties("newman", 80, 2, 1)
	gProps := defineContainerGroup("harness", location, "harnessname", 80, cProps)
	future, err := client.CreateOrUpdate(ctx, resourceGroupName, containerGroupName, gProps)
	errors.HandleError(err)
	err = future.WaitForCompletion(ctx, client.Client)
	errors.HandleError(err)
	result, _ := future.Result(client)
	logger.PrintAndLog(fmt.Sprintf("Container Group: %s has been created.", *result.Name))
}
