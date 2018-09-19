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

func checkStorageAccountNameAvailable(ctx context.Context, storageAccountName *string) bool {
	client := getStorageAccountClient()
	storType := "Microsoft.Storage/storageAccounts"
	checkName := storage.AccountCheckNameAvailabilityParameters{
		Name: storageAccountName,
		Type: &storType,
	}
	data, err := client.CheckNameAvailability(ctx, checkName)
	common.HandleError(err)

	return *data.NameAvailable
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

// CreateStorageAccount - Creates a storage account
func CreateStorageAccount(ctx context.Context, storageAccountName, resourceGroupName, location string) {
	client := getStorageAccountClient()
	nameAvailable := checkStorageAccountNameAvailable(ctx, &storageAccountName)
	if nameAvailable {
		future, err := client.Create(
			ctx,
			resourceGroupName,
			storageAccountName,
			storage.AccountCreateParameters{
				Sku: &storage.Sku{
					Name: storage.StandardLRS},
				Kind:     storage.Storage,
				Location: &location,
				AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{},
			})
		common.HandleError(err)
		err = future.WaitForCompletion(ctx, client.Client)
		common.HandleError(err)
		result, _ := future.Result(client)
		common.PrintAndLog(fmt.Sprintf("Storage Account: %s has been created.", *result.Name))
	} else {
		common.PrintAndLog("Storage Account name is unavailable.")
	}
}
