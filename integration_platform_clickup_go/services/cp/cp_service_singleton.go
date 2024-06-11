package cp_service

import (
	"integration_platform_clickup_go/utils/functions"
	"integration_platform_clickup_go/utils/variables_constant"
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
					"x-api-key":    os.Getenv(variables_constant.CONVISO_PLATFORM_TOKEN_NAME),
				},
				functions.GetFunctionsSingletonInstance(),
			)
		}
	}
	return iCPService
}
