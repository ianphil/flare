package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/iphilpot/flare/apis/storage"
	"github.com/iphilpot/flare/apis/common"
	"github.com/iphilpot/flare/apis/resource"
	"github.com/iphilpot/flare/apis/iam"
	"github.com/spf13/cobra"
)

/*
	Cmd: flare newman --
*/

var postmanCollection string
var postmanIterations int
var storExists bool
var location string

var (
	blobFormatString = `https://%s.blob.core.windows.net`
)

var newman = &cobra.Command{
	Use:   "newman",
	Short: "Run newman collection",
	Long:  "TODO: put more info here",
	Run: func(cmd *cobra.Command, args []string) {
		logger.PrintAndLog("Newman called")
		ctx := context.Background()
		var primaryKey string
		var containerCollection azblob.ContainerURL

		saName, rgName =: common.GenerateNames()

		common.PrintAndLog(fmt.Sprintf("Resource Group: %s | Storage Account: %s", rgName, saName))

		// Create resource group, check if exists
		group := resoure.CreateResourceGroup(ctx, rgName, location)
		
		// Create Storage Account
		storage.CreateStorageAccount(ctx, saName, rgName, location)
		
		// Create storage containers
		storage.CreateStorageContainer(ctx, saName, rgName, "collection")
		storage.CreateStorageContainer(ctx, saName, rgName, "report")


		// first pass primary key from create, if not create need to get for upload
		if primaryKey == "" {
			primaryKey, err = apiStorage.GetStorageAccountPrimaryKey(ctx, &storAccountClient, storName, storName)
		}

		// Upload postman collection file to collection container
		if containerCollection == (azblob.ContainerURL{}) {
			containerCollection = getContainerURL(ctx, storName, "collection", primaryKey)
			fmt.Println(containerCollection.String())
		}

		blobURL := containerCollection.NewBlockBlobURL(postmanCollection)
		file, err := os.Open(postmanCollection)
		if err != nil {
			log.Println(err)
		}

		_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
			BlockSize:   4 * 1024 * 1024,
			Parallelism: 16,
		})

		// All done
		fmt.Println("Completed")

	},
}

func init() {
	rootCmd.AddCommand(newman)
	newman.Flags().StringVarP(&postmanCollection, "collection-file", "c", "", "Postman Collection file name.")
	newman.Flags().StringVarP(&location, "location", "l", "", "Location to put resources. Example: eastus")
	newman.Flags().IntVarP(&postmanIterations, "number", "n", 1, "Number of iterations")
}
