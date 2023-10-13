package Functions

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
	TypeIntegration "integration.platform.clickup/types/type_integration"
	VariablesGlobal "integration.platform.clickup/utils/variables_global"
)

func LoadCustomerByYamlFile() []TypeIntegration.CustomerType {
	// Read the file
	data, err := os.ReadFile("projects.yaml")

	if err != nil {
		fmt.Println("Error ReadFile LoadYamlFileProjects: ", err.Error())
		return nil
	}

	// Create a struct to hold the YAML data
	var projects []TypeIntegration.CustomerType

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(data, &projects)
	if err != nil {
		fmt.Println("Error DataToStruct LoadYamlFileProjects: ", err.Error())
		return nil
	}

	return projects
}

func CustomerExistsYamlFileByClickUpListId(clickUpListId string, customers []TypeIntegration.CustomerType) (result bool) {
	result = false
	for _, customer := range customers {
		if customer.ClickUpListId == clickUpListId {
			result = true
			VariablesGlobal.Customer = customer
			break
		}
	}
	return result
}

func ConvertStringToArrayInt(var1 string) []int {
	var arrayRet []int

	arrayStr := strings.Split(var1, ";")

	for i := 0; i < len(arrayStr); i++ {
		intAux, err := strconv.Atoi(arrayStr[i])
		if err != nil {
			fmt.Println("Error ConvertStringToArrayInt: ", err.Error())
			return nil
		}

		arrayRet = append(arrayRet, intAux)
	}

	return arrayRet
}

func GetTextWithSpace() string {
	ret := ""
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		ret = scanner.Text()
	}
	return ret
}
