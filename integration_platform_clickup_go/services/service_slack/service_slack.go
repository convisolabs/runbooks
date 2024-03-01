package service_slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"integration_platform_clickup_go/types/type_slack"
	"integration_platform_clickup_go/utils/functions"
	"integration_platform_clickup_go/utils/variables_constant"
	"integration_platform_clickup_go/utils/variables_global"
	"io"
	"net/http"
	"os"
)

var globalSlackHeaders map[string]string

func init() {
	globalSlackHeaders = map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + os.Getenv(variables_constant.SLACK_ASSET_TOKEN_NAME),
	}
}

func RequestPostMessage(request type_slack.PostMessage) error {

	var urlPostMessage bytes.Buffer
	urlPostMessage.WriteString(variables_constant.SLACK_API_URL_BASE)
	urlPostMessage.WriteString("chat.postMessage")

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)

	resp, err := functions.HttpRequestRetry(http.MethodPost, urlPostMessage.String(), globalSlackHeaders, payload, *variables_global.Config.ConfclickUp.HttpAttempt)

	if err != nil {
		return errors.New("Error RequestPostMessage: " + err.Error())
	}

	io.ReadAll(resp.Body)

	return nil
}

// func RequestPostMessage(request type_slack.PostMessage) error {
// 	var urlPostMessage bytes.Buffer
// 	urlPostMessage.WriteString(variables_constant.SLACK_API_URL_BASE)
// 	urlPostMessage.WriteString("chat.postMessage")

// 	body, _ := json.Marshal(request)

// 	payload := bytes.NewBuffer(body)

// 	time.Sleep(time.Second)

// 	req, err := http.NewRequest(http.MethodPost, urlPostMessage.String(), payload)
// 	if err != nil {
// 		return errors.New("Error RequestTaskTimeSpent request: " + err.Error())
// 	}

// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+os.Getenv(variables_constant.SLACK_ASSET_TOKEN_NAME))
// 	client := &http.Client{Timeout: time.Second * 10}
// 	resp, err := client.Do(req)
// 	defer req.Body.Close()

// 	if resp.StatusCode != 200 {
// 		return errors.New("Error RequestTaskTimeSpent StatusCode: " + http.StatusText(resp.StatusCode))
// 	}

// 	if err != nil {
// 		return errors.New("Error RequestTaskTimeSpent response: " + err.Error())
// 	}

// 	io.ReadAll(resp.Body)

// 	return nil
// }
