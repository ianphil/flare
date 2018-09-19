package storage

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/2018-03-28/azblob"
	"github.com/iphilpot/flare/apis/errors"
	"github.com/iphilpot/flare/apis/logger"
)

var (
	blobFormatString = `https://%s.blob.core.windows.net`
)

// CreateStorageContainer - Creates storage container
func CreateStorageContainer(ctx context.Context, storageAccountName, resourceGroupName, storageContainerName string) {
	storageContainer := getContainerURL(ctx, storageAccountName, resourceGroupName, storageContainerName)
	_, err := storageContainer.Create(ctx, azblob.Metadata{}, azblob.PublicAccessContainer)
	if err != nil {
		if stErr, ok := err.(azblob.StorageError); ok {
			code := stErr.ServiceCode()
			if code != azblob.ServiceCodeContainerAlreadyExists {
				errors.HandleError(err)
			} else {
				logger.PrintAndLog("Container already exists")
			}
		}
	}
}

func getContainerURL(ctx context.Context, storageAccountName, resourceGroupName, storageContainerName string) azblob.ContainerURL {
	primaryKey := GetStorageAccountPrimaryKey(ctx, storageAccountName, resourceGroupName)
	blobCred, err := azblob.NewSharedKeyCredential(storageAccountName, primaryKey)
	errors.HandleError(err)
	accountURL, _ := url.Parse(fmt.Sprintf(blobFormatString, storageAccountName))
	pipline := azblob.NewPipeline(blobCred, azblob.PipelineOptions{})
	service := azblob.NewServiceURL(*accountURL, pipline)
	return service.NewContainerURL(storageContainerName)
}
