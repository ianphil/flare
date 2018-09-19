package storage

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/iphilpot/flare/apis/config"
	"github.com/iphilpot/flare/apis/errors"
	"github.com/iphilpot/flare/apis/iam"
	"github.com/iphilpot/flare/apis/logger"

	"github.com/Azure/azure-sdk-for-go/services/preview/storage/mgmt/2018-03-01-preview/storage"
	"github.com/Azure/go-autorest/autorest"
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
	errors.HandleError(err)
	if data.Reason != "AlreadyExists" && *data.NameAvailable == false { // This is just in case the name has bad chars
		err := fmt.Errorf("error: Storage Account Name is has invlaid characters")
		errors.HandleError(err)
	}
	return *data.NameAvailable
}

// GetAccountKeys - Returns keys for specified account
func GetAccountKeys(ctx context.Context, storageAccountName, resourceGroupName string) (storage.AccountListKeysResult, error) {
	client := getStorageAccountClient()
	return client.ListKeys(ctx, resourceGroupName, storageAccountName)
}

// GetStorageAccountPrimaryKey - Return primary key of storage account.
func GetStorageAccountPrimaryKey(ctx context.Context, storageAccountName, resourceGroupName string) string {
	keyResult, err := GetAccountKeys(ctx, storageAccountName, resourceGroupName)
	errors.HandleError(err)
	primaryKey := *(((*keyResult.Keys)[0]).Value)
	logger.PrintAndLog("Primary storage account key")
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
		errors.HandleError(err)
		err = future.WaitForCompletion(ctx, client.Client)
		errors.HandleError(err)
		result, _ := future.Result(client)
		logger.PrintAndLog(fmt.Sprintf("Storage Account: %s has been created.", *result.Name))
	} else {
		logger.PrintAndLog("Storage Account has already been created")
	}
}

func logRequest() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				log.Println(err)
			}
			dump, _ := httputil.DumpRequestOut(r, true)
			log.Println(string(dump))
			return r, err
		})
	}
}

func logResponse() autorest.RespondDecorator {
	return func(p autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(r *http.Response) error {
			err := p.Respond(r)
			if err != nil {
				log.Println(err)
			}
			dump, _ := httputil.DumpResponse(r, true)
			log.Println(string(dump))
			return err
		})
	}
}
