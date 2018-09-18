package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// PrintAndLog writes to logger and stdout
func PrintAndLog(message string) {
	log.Println(message)
	fmt.Println(message)
}

// ReadJSON reads a json file, and unmashals it.
// Very useful for template deployments.
func ReadJSON(path string) (*map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read template file: %v\n", err)
	}
	contents := make(map[string]interface{})
	json.Unmarshal(data, &contents)
	return &contents, nil
}

// HandleError - prints and logs error
func HandleError(err error) {
	if err != nil {
		log.Fatalf("failed to read template file: %v\n", err)
	}
}
