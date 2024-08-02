package functions

import (
	"integration_platform_clickup_go/types/type_config"
	"io"
	"net/http"
)

type IFunctions interface {
	LoadConfigsByYamlFile() (type_config.ConfigType, error)
	GetTextWithSpace(label string) string
	HttpRequestRetry(httpMethod string, httpUrl string, headers map[string]string, payload io.Reader, attempt int) (*http.Response, error)
	ConvertStringToArrayInt(var1 string) []int
	Log(text string, onlyScreen bool, saveFile bool) (bool, error)
}
