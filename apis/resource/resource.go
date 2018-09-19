package resource

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/iphilpot/flare/apis/common"
	"github.com/iphilpot/flare/apis/iam"
)

func getResourceAccountClient() resources.GroupsClient {
	client := resources.NewGroupsClient(subID)
	client.Authorizer = iam.GetAuthorizerFromEnvironment()
	return client
}

func checkResourceGroupExists(ctx context.Context, resourceGroupName string) bool {
	client := getResourceAccountClient()
	rgCheck, err := client.CheckExistence(ctx, resourceGroupName)
	common.HandleError(err)
	if rgCheck.StatusCode == 404 {
		return false
	}
	return true
}

// CreateResourceGroup - Creates RG if it doesn't exist
func CreateResourceGroup(ctx context.Context, resourceGroupName, location string) resources.Group {
	client := getResourceAccountClient()
	rgCheck := checkResourceGroupExists(resourceGroupName)
	var group resources.Group
	if rgCheck {
		group, err = client.CreateOrUpdate(ctx, resourceGroupName, resources.Group{Location: &location})
		common.HandleError(err)
	} else {
		group, err = client.Get(ctx, resourceGroupName)
		common.HandleError(err)
	}
	return group
}
