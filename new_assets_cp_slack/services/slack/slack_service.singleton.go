package slack_service

import (
	"new_assets_cp_slack/utils/constants"
	"new_assets_cp_slack/utils/functions"
	"os"
	"sync"
)

var lock = &sync.Mutex{}

var iSlackService ISlackService

func GetSlackServiceSingletonInstance() ISlackService {
	if iSlackService == nil {
		lock.Lock()
		defer lock.Unlock()
		if iSlackService == nil {

			iSlackService = SlackServiceNew(
				map[string]string{
					"Content-Type":  "application/json",
					"Authorization": "Bearer " + os.Getenv(constants.SLACK_ASSET_TOKEN_NAME),
				},
				functions.GetFunctionsSingletonInstance(),
			)
		}
	}
	return iSlackService
}
