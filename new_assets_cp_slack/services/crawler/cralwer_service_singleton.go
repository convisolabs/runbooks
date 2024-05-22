package crawler_service

import (
	asset_repository "new_assets_cp_slack/repositories/assets"
	slack_service "new_assets_cp_slack/services/slack"
	"sync"
)

var lock = &sync.Mutex{}

var iCrawlerService ICrawlerService

func GetCrawlerServiceSingletonInstance() ICrawlerService {
	if iCrawlerService == nil {
		lock.Lock()
		defer lock.Unlock()
		if iCrawlerService == nil {

			iCrawlerService = CrawlerServiceNew(
				slack_service.GetSlackServiceSingletonInstance(),
				asset_repository.GetAssetRepositorySingletonInstance(),
			)
		}
	}
	return iCrawlerService
}
