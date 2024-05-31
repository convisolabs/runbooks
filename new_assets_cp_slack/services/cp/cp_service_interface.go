package cp_service

import type_cp "new_assets_cp_slack/types/cp"

type ICPService interface {
	GetAssetsByTime(parameter type_cp.AssetsByTimeParameters) ([]type_cp.Asset, error)
}
