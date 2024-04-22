package variables_global

import (
	"integration_platform_clickup_go/types/type_config"
)

var Customer = type_config.ConfigTypeIntegration{
	IntegrationName:     "",
	PlatformID:          0,
	ClickUpListId:       "",
	ClickUpCustomerList: "No selected project",
	//CheckTagsValidationStory: "",
}

var Config type_config.ConfigType
