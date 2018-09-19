package storage

import (
	"context"
	"os"

	"github.com/Azure/azure-storage-blob-go/2018-03-28/azblob"
	"github.com/iphilpot/flare/apis/errors"
)

// UploadBlob - uploads a file to Azure Storage
func UploadBlob(ctx context.Context, storageAccountName, resourceGroupName, storageCollectionName, postmanCollectionFile string) {
	containerCollection := getContainerURL(ctx, storageAccountName, resourceGroupName, storageCollectionName)
	blobURL := containerCollection.NewBlockBlobURL(postmanCollectionFile)
	file, err := os.Open(postmanCollectionFile)
	errors.HandleError(err)
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16,
	})
	errors.HandleError(err)
}
