package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/storage/mgmt/2018-03-01-preview/storage"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/spf13/cobra"
)

/*
	Cmd: flare newman --
*/

var postmanCollection string
var postmanIterations int
var storExists bool
var location string

var newman = &cobra.Command{
	Use:   "newman",
	Short: "Run newman collection",
	Long:  "TODO: put more info here",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Newman called")
		ctx := context.Background()

		// create an authorizer from env vars or Azure Managed Service Identity
		authorizer, err := auth.NewAuthorizerFromEnvironment()
		if err != nil {
			fmt.Println("Not Authorized")
		}

		// Get envvars and generate storage account name based on last part of subID
		subID := os.Getenv("AZURE_SUB_ID")
		storName := strings.ToLower(fmt.Sprintf("flare%s", strings.Split(subID, "-")[4]))
		println(storName)

		// Create resource group, check if exists
		groupsClient := resources.NewGroupsClient(subID)
		groupsClient.Authorizer = authorizer

		rgCheck, err := groupsClient.CheckExistence(ctx, storName)
		if err != nil {
			log.Println(err)
		}

		var group resources.Group
		if rgCheck.StatusCode == 404 {
			group, err = groupsClient.CreateOrUpdate(ctx, storName, resources.Group{Location: &location})
			if err != nil {
				log.Println(err)
			}
		} else {
			group, err = groupsClient.Get(ctx, storName)
			if err != nil {
				log.Println(err)
			}
		}

		fmt.Println(*group.Name)

		// Test if storage exists
		storAccountClient := storage.NewAccountsClient(subID)
		storAccountClient.Authorizer = authorizer
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
					Kind:                              storage.Storage,
					Location:                          &location,
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
		}

		// Get Primary storage account key
		keyResult, err := storAccountClient.ListKeys(ctx, storName, storName)
		if err != nil {
			log.Println(err)
		}

		primaryKey := *(((*keyResult.Keys)[0]).Value)

		fmt.Printf("Primary storage account key: %s\n", primaryKey)

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
