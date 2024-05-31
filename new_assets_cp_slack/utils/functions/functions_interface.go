package functions

import (
	"io"
	"net/http"
	type_config "new_assets_cp_slack/types/config"
	type_integration "new_assets_cp_slack/types/integration"
)

type IFunctions interface {
	SaveYamlFile(params type_integration.SaveFile) error
	LoadConfigsByYamlFile() (type_config.ConfigType, error)
	HttpRequestRetry(httpMethod string, httpUrl string, headers map[string]string, payload io.Reader, attempt int) (*http.Response, error)
}
