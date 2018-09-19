package resource

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/iphilpot/flare/apis/config"
	"github.com/iphilpot/flare/apis/errors"
	"github.com/iphilpot/flare/apis/iam"
)

func getResourceAccountClient() resources.GroupsClient {
	c := config.GetConfig()
	client := resources.NewGroupsClient(c.AzureSubscriptionID)
	client.Authorizer = iam.GetAuthorizerFromEnvironment()
	return client
}

func checkResourceGroupExists(ctx context.Context, resourceGroupName string) bool {
	client := getResourceAccountClient()
	rgCheck, err := client.CheckExistence(ctx, resourceGroupName)
	errors.HandleError(err)
	if rgCheck.StatusCode == 404 {
		return false
	}
	return true
}

// CreateResourceGroup - Creates RG if it doesn't exist
func CreateResourceGroup(ctx context.Context, resourceGroupName, location string) resources.Group {
	client := getResourceAccountClient()
	rgCheck := checkResourceGroupExists(ctx, resourceGroupName)
	var group resources.Group
	var err error
	if !rgCheck {
		group, err = client.CreateOrUpdate(ctx, resourceGroupName, resources.Group{Location: &location})
		errors.HandleError(err)
	} else {
		group, err = client.Get(ctx, resourceGroupName)
		errors.HandleError(err)
	}
	return group
}
