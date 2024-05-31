package cp_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	type_cp "new_assets_cp_slack/types/cp"
	"new_assets_cp_slack/utils/constants"
	"new_assets_cp_slack/utils/functions"
)

const CONVISO_PLATFORM_ASSETS_BY_TIME = `
	query Assets($CompanyId:ID!,$Page:Int,$Limit:Int,$Search:AssetsSearch)
	{
		assets(companyId: $CompanyId, page: $Page, limit: $Limit, search: $Search) {
			collection {
				id
				name
				createdAt
			}
			metadata {
				currentPage
				limitValue
				totalCount
				totalPages
			}
		}
	}
`

type CPService struct {
	HttpHeaders map[string]string
	functions   functions.IFunctions
}

func CPServiceNew(HttpHeaders map[string]string, functions functions.IFunctions) ICPService {
	return &CPService{
		HttpHeaders: HttpHeaders,
		functions:   functions,
	}
}

func (f *CPService) GetAssetsByTime(parameter type_cp.AssetsByTimeParameters) ([]type_cp.Asset, error) {

	var result type_cp.AssetsResponse
	var ret []type_cp.Asset

	for i := 0; i <= result.Data.Assets.Metadata.TotalPages; i++ {
		parameter.Page = parameter.Page + 1
		parameters, _ := json.Marshal(parameter)
		body, _ := json.Marshal(map[string]string{
			"query":     CONVISO_PLATFORM_ASSETS_BY_TIME,
			"variables": string(parameters),
		})

		payload := bytes.NewBuffer(body)

		resp, err := f.functions.HttpRequestRetry(http.MethodPost, constants.CONVISO_PLATFORM_API_GRAPHQL, f.HttpHeaders, payload, 3)

		if err != nil {
			fmt.Println("Error GetAssetsByTime HttpRequestRetry: ", err.Error())
			return ret, err
		}

		data, _ := io.ReadAll(resp.Body)

		json.Unmarshal([]byte(string(data)), &result)

		ret = append(ret, result.Data.Assets.Collection...)

	}

	return ret, nil
}
