package clickup_service

import (
	cp_service "integration_platform_clickup_go/services/cp"
	"integration_platform_clickup_go/utils/functions"
	"integration_platform_clickup_go/utils/variables_constant"
	"os"
	"sync"
)

var lock = &sync.Mutex{}

var iClickupService IClickupService

func GetClickupServiceSingletonInstance() IClickupService {
	if iClickupService == nil {
		lock.Lock()
		defer lock.Unlock()
		if iClickupService == nil {

			iClickupService = ClickupServiceNew(
				map[string]string{
					"Content-Type":  "application/json",
					"Authorization": os.Getenv(variables_constant.CLICKUP_TOKEN_NAME),
				},
				functions.GetFunctionsSingletonInstance(),
				cp_service.GetCPServiceSingletonInstance(),
			)
		}
	}
	return iClickupService
}
