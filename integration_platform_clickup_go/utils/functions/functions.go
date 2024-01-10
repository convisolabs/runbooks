package functions

import (
	"bufio"
	"errors"
	"fmt"
	"integration_platform_clickup_go/types/type_config"
	"integration_platform_clickup_go/utils/variables_global"
	"os"
	"runtime"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

func LoadConfigsByYamlFile() (type_config.ConfigType, error) {

	// Create a struct to hold the YAML data
	var config type_config.ConfigType

	// Read the file
	data, err := os.ReadFile("projects.yaml")

	if err != nil {
		return config, errors.New("Error ReadFile " + err.Error())
	}

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, errors.New("Error DataToStruct " + err.Error())
	}

	return config, nil
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

// func GetTextWithSpace() string {
// 	ret := ""
// 	scanner := bufio.NewScanner(os.Stdin)
// 	if scanner.Scan() {
// 		ret = scanner.Text()
// 	}
// 	return ret
// }

func GetTextWithSpace(label string) string {
	ret := ""

	EOL := byte('\r')

	if strings.ToLower(runtime.GOOS) == "linux" {
		EOL = byte('\n')
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(label)

	ret, error := reader.ReadString(EOL)

	if error != nil {
		fmt.Print("Error function GetTextWithSpace ", error)
		return ret
	}

	ret = strings.Replace(ret, "\r", "", -1)
	ret = strings.Replace(ret, "\n", "", -1)

	return ret
}
