package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/iphilpot/flare/apis/config"
	"github.com/iphilpot/flare/apis/logger"
)

// PrintAndLog writes to logger and stdout
func PrintAndLog(message string) {
	logger.Log.Println(message)
	fmt.Println(message)
}

// ReadJSON reads a json file, and unmashals it.
// Very useful for template deployments.
func ReadJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	HandleError(err)
	contents := make(map[string]interface{})
	json.Unmarshal(data, &contents)
	return &contents, nil
}

// HandleError - prints and logs error
func HandleError(err error) {
	if err != nil {
		logger.Log.Fatalln(err)
	}
}

// GenerateNames - Generates and returns SA Name and RG Name
func GenerateNames() (storageAccountName, resourceGroupName string) {
	subscriptionIDLastPart := strings.Split(config.Config.SubscriptionID, "-")[4]
	storageAccountName := strings.ToLower(fmt.Sprintf("flare-sa-%s", subscriptionIDLastPart))
	resourceGroupName := strings.ToLower(fmt.Sprintf("flare-rg-%s", subscriptionIDLastPart))
	return storageAccountName, resourceGroupName
}
