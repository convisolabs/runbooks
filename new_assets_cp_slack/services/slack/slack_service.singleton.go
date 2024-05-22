package slack_service

import "sync"

var lock = &sync.Mutex{}

var iSlackService ISlackService

func GetSlackServiceSingletonInstance() ISlackService {
	if iSlackService == nil {
		lock.Lock()
		defer lock.Unlock()
		if iSlackService == nil {

			iSlackService = SlackServiceNew()
		}
	}
	return iSlackService
}
