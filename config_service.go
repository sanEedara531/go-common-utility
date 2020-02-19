package common

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"os"
)

// ConfigurationObj exported
var ConfigurationObj Configurations

// TableDetails exported
var TableDetails TableConfigurations

// ConfigObj exported
func ConfigObj(filePath string) Configurations {
	var configurations Configurations

	env := os.Getenv("env")

	//set default environment variable
	if env == "" {
		env = "dev"
	}

	jsonFile, err := os.Open(filePath + "config/config." + env + ".json")

	if err != nil {
		fmt.Println(err)
	}

	//close json file after main file executes completely
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &configurations)
	//	if err != nil {  return err }

	return configurations
}