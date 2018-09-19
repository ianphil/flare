package storage

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/2018-03-28/azblob"
	"github.com/iphilpot/flare/apis/common"
)

var (
	blobFormatString = `https://%s.blob.core.windows.net`
)

// CreateStorageContainer - Creates storage container
func CreateStorageContainer(ctx context.Context, storageAccountName, resourceGroupName, storageContainerName string) {
	primaryKey := GetStorageAccountPrimaryKey(ctx, storageAccountName, resourceGroupName)
	storageContainer = getContainerURL(ctx, storageAccountName, storageContainerName, primaryKey)
	_, err = containerCollection.Create(ctx, azblob.Metadata{}, azblob.PublicAccessContainer)
	common.HandleError(err)
}

func getContainerURL(ctx context.Context, storageAccountName, storageContainerName, primaryKey string) azblob.ContainerURL {
	blobCred, err := azblob.NewSharedKeyCredential(accountName, primaryKey)
	common.HandleError(err)
	accountURL, _ := url.Parse(fmt.Sprintf(blobFormatString, storageAccountName))
	pipline := azblob.NewPipeline(blobCred, azblob.PipelineOptions{})
	service := azblob.NewServiceURL(*accountURL, pipline)
	return service.NewContainerURL(storageContainerName)
}
