package functions

import (
	"bufio"
	"errors"
	"fmt"
	"integration_platform_clickup_go/types/type_config"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Functions struct{}

func FunctionsNew() IFunctions {
	return &Functions{}
}

func (f *Functions) LoadConfigsByYamlFile() (type_config.ConfigType, error) {

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

// func CustomerExistsYamlFileByClickUpListId(clickUpListId string, customers []type_config.ConfigTypeIntegration) (result bool) {
// 	result = false
// 	for _, customer := range customers {
// 		if customer.ClickUpListId == clickUpListId {
// 			result = true
// 			variables_global.Customer = customer
// 			break
// 		}
// 	}
// 	return result
// }

func (f *Functions) ConvertStringToArrayInt(var1 string) []int {
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

func (f *Functions) GetTextWithSpace(label string) string {
	ret := ""

	EOL := byte('\r')

	if !strings.EqualFold(runtime.GOOS, "windows") {
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

func (f *Functions) HttpRequestRetry(httpMethod string, httpUrl string, headers map[string]string, payload io.Reader, attempt int) (*http.Response, error) {
	req, err := http.NewRequest(httpMethod, httpUrl, payload)

	msgError := ""

	if err != nil {
		return nil, errors.New("Error HttpRequestRetry Request: " + err.Error())
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	aux := 1
	for ok := true; ok; ok = (attempt + 1) > aux {
		aux = aux + 1

		client := &http.Client{Timeout: time.Second * 10}
		resp, err := client.Do(req)

		if err != nil {
			return resp, errors.New("Error HttpRequestRetry ClientDo: " + err.Error())
		}

		if resp.StatusCode != 200 {
			time.Sleep(time.Second)
			msgError = msgError + "Retry (" + string(aux) + "): " + http.StatusText(resp.StatusCode) + " "
			continue
		}

		return resp, nil
	}

	return nil, errors.New("Error HttpRequestRetry Final: " + msgError)
}

func (f *Functions) WriteFile(fileName string, line string) (bool, error) {

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		file, err = os.Create(fileName)
		if err != nil {
			fmt.Println(err.Error())
			return false, errors.New("Error WriteFile - Create: " + err.Error())
		}
	}

	_, err = file.WriteString(line + "\n")

	if err != nil {
		fmt.Println(err.Error())
		defer file.Close()
		return false, errors.New("Error WriteFile - WriteString: " + err.Error())
	}

	defer file.Close()

	return true, nil
}

func (f *Functions) Log(text string, onlyScreen bool, saveFile bool) (bool, error) {

	if !onlyScreen && saveFile {
		fileName := "integration_platform_clickup_go_" + time.Now().Format("20060102")
		_, err := f.WriteFile(fileName, text)
		if err != nil {
			return false, errors.New("Log - " + err.Error())
		}
	}

	fmt.Println(time.Now().Format("2006-01-02 15:04:05")+": ", text)

	return true, nil
}
