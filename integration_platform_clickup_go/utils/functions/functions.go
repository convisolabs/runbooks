package functions

import (
	"bufio"
	"fmt"
	"integration_platform_clickup_go/types/type_config"
	"integration_platform_clickup_go/utils/variables_global"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

func LoadConfigsByYamlFile() type_config.ConfigType {

	// Create a struct to hold the YAML data
	var config type_config.ConfigType

	// Read the file
	data, err := os.ReadFile("projects.yaml")

	if err != nil {
		fmt.Println("Error ReadFile LoadConfigsByYamlFile: ", err.Error())
		return config
	}

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error DataToStruct LoadConfigsByYamlFile: ", err.Error())
		return config
	}

	return config
}

func CustomerExistsYamlFileByClickUpListId(clickUpListId string, customers []type_config.ConfigTypeIntegration) (result bool) {
	result = false
	for _, customer := range customers {
		if customer.ClickUpListId == clickUpListId {
			result = true
			variables_global.Customer = customer
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
