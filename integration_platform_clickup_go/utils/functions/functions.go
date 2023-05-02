package functions

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	VariablesGlobal "integration.platform.clickup/utils/variables_global"
)

func LoadCustomerByYamlFile() []VariablesGlobal.CustomerType {
	// Read the file
	data, err := ioutil.ReadFile("projects.yaml")
	if err != nil {
		fmt.Println("Error ReadFile LoadYamlFileProjects: ", err.Error())
		return nil
	}

	// Create a struct to hold the YAML data
	var projects []VariablesGlobal.CustomerType

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(data, &projects)
	if err != nil {
		fmt.Println("Error DataToStruct LoadYamlFileProjects: ", err.Error())
		return nil
	}

	return projects
}

func CustomerExistsYamlFileByClickUpListId(clickUpListId string, customers []VariablesGlobal.CustomerType) (result bool) {
	result = false
	for _, customer := range customers {
		if customer.ClickUpListId == clickUpListId {
			result = true
			break
		}
	}
	return result
}
