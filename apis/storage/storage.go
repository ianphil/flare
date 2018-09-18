package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/storage/mgmt/2018-03-01-preview/storage"
)

func getStorageAccountClient() storage.AccountsClient {
	storageAccountClient := storage.NewAccountsClient("")

	return storageAccountClient
}

// GetStorageAccountPrimaryKey - Return primary key of storage account.
func GetStorageAccountPrimaryKey(ctx context.Context, storageAccountClient *storage.AccountsClient, accountName, accountGroupName string) (string, error) {
	// Get Primary storage account key
	keyResult, err := storageAccountClient.ListKeys(ctx, accountGroupName, accountName)
	if err != nil {
		log.Println(err)
	}

	primaryKey := *(((*keyResult.Keys)[0]).Value)

	fmt.Printf("Primary storage account key: %s\n", primaryKey)

	return primaryKey, err
}
