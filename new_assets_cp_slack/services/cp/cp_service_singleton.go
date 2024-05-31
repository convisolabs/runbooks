package cp_service

import (
	"new_assets_cp_slack/utils/constants"
	"new_assets_cp_slack/utils/functions"
	"os"
	"sync"
)

var lock = &sync.Mutex{}

var iCPService ICPService

func GetCPServiceSingletonInstance() ICPService {
	if iCPService == nil {
		lock.Lock()
		defer lock.Unlock()
		if iCPService == nil {

			iCPService = CPServiceNew(
				map[string]string{
					"Content-Type": "application/json",
					"x-api-key":    os.Getenv(constants.CONVISO_PLATFORM_TOKEN_NAME),
				},
				functions.GetFunctionsSingletonInstance(),
			)
		}
	}
	return iCPService
}
