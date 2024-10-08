package slack_service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	type_slack "new_assets_cp_slack/types/slack"
	"new_assets_cp_slack/utils/constants"
	"new_assets_cp_slack/utils/functions"
	"new_assets_cp_slack/utils/globals"
)

type SlackService struct {
	functions   functions.IFunctions
	HttpHeaders map[string]string
}

func SlackServiceNew(HttpHeaders map[string]string, functions functions.IFunctions) ISlackService {

	return &SlackService{
		functions:   functions,
		HttpHeaders: HttpHeaders,
	}
}

func (f *SlackService) RequestPostMessage(request type_slack.PostMessage) error {

	var urlPostMessage bytes.Buffer
	urlPostMessage.WriteString(constants.SLACK_API_URL_BASE)
	urlPostMessage.WriteString("chat.postMessage")

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)

	resp, err := f.functions.HttpRequestRetry(http.MethodPost, urlPostMessage.String(), f.HttpHeaders, payload, *globals.Config.ConfigSlack.HttpAttempt)

	if err != nil {
		return errors.New("Error RequestPostMessage: " + err.Error())
	}

	io.ReadAll(resp.Body)

	return nil
}
