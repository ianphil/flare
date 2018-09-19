package iam

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/iphilpot/flare/apis/errors"
)

// GetAuthorizerFromEnvironment - returns authorizer from Env Vars
func GetAuthorizerFromEnvironment() autorest.Authorizer {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	errors.HandleError(err)

	return authorizer
}
