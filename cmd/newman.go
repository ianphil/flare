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

	"github.com/Azure/azure-sdk-for-go/services/preview/storage/mgmt/2018-03-01-preview/storage"
	"github.com/Azure/azure-storage-blob-go/2018-03-28/azblob"
	"github.com/Azure/go-autorest/autorest"
	"github.com/iphilpot/flare/apis/common"
	"github.com/iphilpot/flare/apis/resource"
	"github.com/iphilpot/flare/apis/iam"
	"github.com/spf13/cobra"

	apiStorage "github.com/iphilpot/flare/apis/storage"
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
		common.PrintAndLog("Newman called")
		ctx := context.Background()
		var primaryKey string
		var containerCollection azblob.ContainerURL

		saName, rgName =: common.GenerateNames()

		common.PrintAndLog(fmt.Sprintf("Resource Group: %s | Storage Account: %s", rgName, saName))

		// Create resource group, check if exists
		group := resoure.CreateResourceGroup(ctx, rgName, location)

		common.PrintAndLog(fmt.Sprintf("Resource Group %s created", *group.Name))

		// Test if storage exists
		storAccountClient := storage.NewAccountsClient(subID)
		storAccountClient.Authorizer = iam.GetAuthorizerFromEnvironment()
		//storAccountClient.RequestInspector = logRequest()
		//storAccountClient.ResponseInspector = logResponse()

		storType := "Microsoft.Storage/storageAccounts"
		checkName := storage.AccountCheckNameAvailabilityParameters{
			Name: &storName,
			Type: &storType,
		}

		// Check name, I will assume if it's false that the name exists in current sub
		// probably should harden this logic but it's almost safe because of the name/sub
		// relationship.
		data, err := storAccountClient.CheckNameAvailability(ctx, checkName)
		if err != nil {
			log.Println(err)
		}

		// Create Storage Account
		if *data.NameAvailable {
			fmt.Println("Storeage name is available")

			future, err := storAccountClient.Create(
				ctx,
				storName,
				storName,
				storage.AccountCreateParameters{
					Sku: &storage.Sku{
						Name: storage.StandardLRS},
					Kind:     storage.Storage,
					Location: &location,
					AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{},
				})
			if err != nil {
				log.Println(err)
			}

			err = future.WaitForCompletion(ctx, storAccountClient.Client)
			if err != nil {
				log.Println(err)
			}

			result, _ := future.Result(storAccountClient)

			fmt.Println(*result.Name)

			// Get Primary Key
			primaryKey, err = apiStorage.GetStorageAccountPrimaryKey(ctx, &storAccountClient, storName, storName)

			// Create storage containers
			containerCollection = getContainerURL(ctx, storName, "collection", primaryKey)
			containerReport := getContainerURL(ctx, storName, "report", primaryKey)

			_, err = containerCollection.Create(ctx, azblob.Metadata{}, azblob.PublicAccessContainer)
			if err != nil {
				log.Println(err)
			}

			_, err = containerReport.Create(ctx, azblob.Metadata{}, azblob.PublicAccessContainer)
			if err != nil {
				log.Println(err)
			}
		}

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

func getContainerURL(ctx context.Context, accountName, containerName, primaryKey string) azblob.ContainerURL {
	blobCred, err := azblob.NewSharedKeyCredential(accountName, primaryKey)
	if err != nil {
		log.Println(err)
	}

	accountURL, _ := url.Parse(fmt.Sprintf(blobFormatString, accountName))
	pipline := azblob.NewPipeline(blobCred, azblob.PipelineOptions{})
	service := azblob.NewServiceURL(*accountURL, pipline)
	return service.NewContainerURL(containerName)
}

func init() {
	rootCmd.AddCommand(newman)
	newman.Flags().StringVarP(&postmanCollection, "collection-file", "c", "", "Postman Collection file name.")
	newman.Flags().StringVarP(&location, "location", "l", "", "Location to put resources. Example: eastus")
	newman.Flags().IntVarP(&postmanIterations, "number", "n", 1, "Number of iterations")
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
