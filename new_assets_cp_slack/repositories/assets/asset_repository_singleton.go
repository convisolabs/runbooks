package asset_repository

import "sync"

var lock = &sync.Mutex{}

var iAssetReposiroty IAssetReposiroty

func GetAssetRepositorySingletonInstance() IAssetReposiroty {
	if iAssetReposiroty == nil {
		lock.Lock()
		defer lock.Unlock()
		if iAssetReposiroty == nil {

			iAssetReposiroty = AssetRepositoryNew()
		}
	}
	return iAssetReposiroty
}
