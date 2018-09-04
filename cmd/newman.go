package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

/*
	Cmd: flare newman --
*/

var postmanCollection string
var postmanIterations int

var newman = &cobra.Command{
	Use:   "newman",
	Short: "Run newman collection",
	Long:  "TODO: put more info here",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Newman called")

		// Get envvars
		subID := os.Getenv("AZURE_SUB_ID")
		storName := fmt.Sprintf("flare-%s", strings.Split(subID, "-")[4])
		println(storName)
	},
}

func init() {
	rootCmd.AddCommand(newman)
	newman.Flags().StringVarP(&postmanCollection, "collection-file", "c", "", "Postman Collection file name.")
	newman.Flags().IntVarP(&postmanIterations, "number", "n", 1, "Number of iterations")
}
