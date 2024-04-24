package slack_service

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

type SlackService struct{}

var globalSlackHeaders map[string]string

func SlackServiceNew() ISlackService {
	globalSlackHeaders = map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + os.Getenv(variables_constant.SLACK_ASSET_TOKEN_NAME),
	}

	return &SlackService{}
}

func (s *SlackService) RequestPostMessage(request type_slack.PostMessage) error {

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
