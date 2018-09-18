package config

import (
	"github.com/gobuffalo/envy"
	"github.com/iphilpot/flare/apis/common"
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
	common.HandleError(err)

	azClientID, err := envy.MustGet("AZURE_CLIENT_ID")
	common.HandleError(err)

	azClientSecret, err := envy.MustGet("AZURE_CLIENT_SECRET")
	common.HandleError(err)

	azSubscriptionID, err := envy.MustGet("AZURE_SUB_ID")
	common.HandleError(err)

	return Config{
		AzureTenantID:       azTenantID,
		AzureClientID:       azClientID,
		AzureClientSecret:   azClientSecret,
		AzureSubscriptionID: azSubscriptionID,
	}
}
