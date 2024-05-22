package asset_repository

import type_repository "new_assets_cp_slack/types/repository"

type IAssetReposiroty interface {
	Insert(asset type_repository.Asset) error
	AssetExist(asset type_repository.Asset) (bool, error)
}
