package cmd

import (
	"context"
	"fmt"

	"github.com/iphilpot/flare/apis/common"
	"github.com/iphilpot/flare/apis/containers"
	"github.com/iphilpot/flare/apis/logger"
	"github.com/iphilpot/flare/apis/resource"
	"github.com/iphilpot/flare/apis/storage"
	"github.com/spf13/cobra"
)

/*
	Cmd: flare newman --
*/

var (
	blobFormatString  = `https://%s.blob.core.windows.net`
	postmanCollection string
	postmanIterations int
	storExists        bool
	location          string
)

var newman = &cobra.Command{
	Use:   "newman",
	Short: "Run newman collection",
	Long:  "TODO: put more info here",
	Run: func(cmd *cobra.Command, args []string) {
		logger.PrintAndLog("Newman called")
		ctx := context.Background()

		saName, rgName, dnsName := common.GenerateNames()

		logger.PrintAndLog(fmt.Sprintf("Resource Group: %s | Storage Account: %s", rgName, saName))

		// Create resource group, check if exists
		_ = resource.CreateResourceGroup(ctx, rgName, location)

		// Create Storage Account
		storage.CreateStorageAccount(ctx, saName, rgName, location)

		// Create storage containers
		storage.CreateStorageContainer(ctx, saName, rgName, "collection")
		storage.CreateStorageContainer(ctx, saName, rgName, "report")

		// Upload postman collection
		storage.UploadBlob(ctx, saName, rgName, "collection", postmanCollection)

		// Create Newman container in ACI
		containerGroup := containers.ContainerGroup{
			GroupName:         "NewmanHarness",
			DNSName:           dnsName,
			Location:          location,
			Port:              80,
			ResourceGroupName: rgName,
			Containers: []containers.Container{
				containers.Container{
					Name:      "newman",
					ImageName: "nginx",
					CPU:       1,
					Memory:    2,
					Port:      8080, // TODO: this is not needed, going to need to create specific container types or function... hmmm
				},
				containers.Container{
					Name:      "web",
					ImageName: "nginx",
					CPU:       1,
					Memory:    2,
					Port:      80,
				},
			},
		}
		containerGroup.CreateContainerGroup(ctx)

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
