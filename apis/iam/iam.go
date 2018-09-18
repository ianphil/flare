package iam

import (
	"github.com/iphilpot/flare/apis/common"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// GetAuthorizerFromEnvironment - returns authorizer from Env Vars
func GetAuthorizerFromEnvironment() autorest.Authorizer {
	// create an authorizer from env vars or Azure Managed Service Identity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		common.PrintAndLog("GetAuthorizerFromEnvironment - Not Authorized")
	}

	return authorizer
}
