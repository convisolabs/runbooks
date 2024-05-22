package functions

import (
	"errors"
	"io"
	"net/http"
	type_config "new_assets_cp_slack/types/config"
	"os"
	"time"

	"gopkg.in/yaml.v3"
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

func HttpRequestRetry(httpMethod string, httpUrl string, headers map[string]string, payload io.Reader, attempt int) (*http.Response, error) {
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
