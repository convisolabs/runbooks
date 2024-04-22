package type_config

type ConfigTypeIntegration struct {
	IntegrationName               string `yaml:"integrationname"`
	PlatformID                    int    `yaml:"platformid"`
	ClickUpListId                 string `yaml:"clickuplistid"`
	ClickUpCustomerList           string `yaml:"clickupcustomerlist"`
	AssetNewFortify               bool   `yaml:"assetNewFortify"`
	ValidateTag                   bool   `yaml:"validateTag"`
	ValidatePSCustomer            bool   `yaml:"validatePSCustomer"`
	ValidatePSConvisoPlatformLink bool   `yaml:"validatePSConvisoPlatformLink"`
}

type ConfigTypeClickup struct {
	User        int64 `yaml:"user"`
	HttpAttempt *int  `yaml:"httpattempt",omitempty`
}

type ConfigTypeGeneral struct {
	IntegrationDefault int `yaml:"integrationdefault"`
}

type ConfigType struct {
	ConfigGeneral ConfigTypeGeneral       `yaml:"configGeneral"`
	ConfclickUp   ConfigTypeClickup       `yaml:"configclickUp"`
	Integrations  []ConfigTypeIntegration `yaml:"integrations"`
	Tags          []ConfigTag             `yaml:"tags"`
}

type ConfigTag struct {
	Value          string `yaml:"value"`
	DeliveryPoints int    `yaml:"deliveryPoints"`
}
