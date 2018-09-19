package common

import (
	"fmt"
	"strings"

	"github.com/iphilpot/flare/apis/config"
)

// GenerateNames - Generates and returns SA Name and RG Name
func GenerateNames() (storageAccountName, resourceGroupName string) {
	c := config.GetConfig()
	subscriptionIDLastPart := strings.Split(c.AzureSubscriptionID, "-")[4]
	storageAccountName = strings.ToLower(fmt.Sprintf("flaresa%s", subscriptionIDLastPart))
	resourceGroupName = strings.ToLower(fmt.Sprintf("flarerg%s", subscriptionIDLastPart))
	return storageAccountName, resourceGroupName
}
