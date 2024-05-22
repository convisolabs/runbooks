package asset_repository

import (
	"database/sql"
	"errors"
	type_repository "new_assets_cp_slack/types/repository"
	"new_assets_cp_slack/utils/constants"
)

type AssetRepository struct{}

func AssetRepositoryNew() IAssetReposiroty {
	return &AssetRepository{}
}

func (f *AssetRepository) Insert(asset type_repository.Asset) error {
	db, err := sql.Open("sqlite3", constants.DATABASE_SQLITE_PATH)

	if err != nil {
		return err
	}

	// insert
	stmt, err := db.Prepare("INSERT INTO assets(id, cp_company_id, cp_asset_id, name, dt_created) values(?,?,?,?,?)")

	if err != nil {
		db.Close()
		return err
	}

	res, err := stmt.Exec(asset.Id, asset.CPCompanyId, asset.CPAssetId, asset.Name, asset.DtCreated)

	if err != nil {
		db.Close()
		return err
	}

	affect, err := res.RowsAffected()

	if err != nil || affect == 0 {
		db.Close()
		return errors.New("Error Asset Insert " + err.Error())
	}

	db.Close()
	return nil
}

func (f *AssetRepository) AssetExist(asset type_repository.Asset) (bool, error) {
	var exist bool

	db, err := sql.Open("sqlite3", constants.DATABASE_SQLITE_PATH)
	if err != nil {
		return false, err
	}

	db.QueryRow("select count(*) > 0 from assets where name=? and cp_company_id=?", asset.Name, asset.CPCompanyId).Scan(&exist)

	db.Close()

	return exist, nil
}
