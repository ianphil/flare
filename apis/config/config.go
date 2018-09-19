package config

import (
	"github.com/gobuffalo/envy"
	"github.com/iphilpot/flare/apis/errors"
)

// Config - containes all env vars
type Config struct {
	AzureTenantID       string
	AzureClientID       string
	AzureClientSecret   string
	AzureSubscriptionID string
}

// GetConfig - Creates new config struct with state
func GetConfig() Config {
	azTenantID, err := envy.MustGet("AZURE_TENANT_ID")
	errors.HandleError(err)

	azClientID, err := envy.MustGet("AZURE_CLIENT_ID")
	errors.HandleError(err)

	azClientSecret, err := envy.MustGet("AZURE_CLIENT_SECRET")
	errors.HandleError(err)

	azSubscriptionID, err := envy.MustGet("AZURE_SUB_ID")
	errors.HandleError(err)

	return Config{
		AzureTenantID:       azTenantID,
		AzureClientID:       azClientID,
		AzureClientSecret:   azClientSecret,
		AzureSubscriptionID: azSubscriptionID,
	}
}
