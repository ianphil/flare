package storage

import (
	"context"
	"fmt"

	"github.com/iphilpot/flare/apis/common"
	"github.com/iphilpot/flare/apis/config"
	"github.com/iphilpot/flare/apis/iam"

	"github.com/Azure/azure-sdk-for-go/services/preview/storage/mgmt/2018-03-01-preview/storage"
)

func getStorageAccountClient() storage.AccountsClient {
	configInfo := config.GetConfig()
	storageAccountClient := storage.NewAccountsClient(configInfo.AzureSubscriptionID)
	storageAccountClient.Authorizer = iam.GetAuthorizerFromEnvironment()
	return storageAccountClient
}

// GetAccountKeys - Returns keys for specified account
func GetAccountKeys(ctx context.Context, accountName, accountGroupName string) (storage.AccountListKeysResult, error) {
	client := getStorageAccountClient()
	return client.ListKeys(ctx, accountGroupName, accountName)
}

// GetStorageAccountPrimaryKey - Return primary key of storage account.
func GetStorageAccountPrimaryKey(ctx context.Context, accountName, accountGroupName string) string {
	keyResult, err := GetAccountKeys(ctx, accountGroupName, accountName)
	common.HandleError(err)
	primaryKey := *(((*keyResult.Keys)[0]).Value)
	common.PrintAndLog(fmt.Sprintf("Primary storage account key: %s\n", primaryKey))
	return primaryKey
}
