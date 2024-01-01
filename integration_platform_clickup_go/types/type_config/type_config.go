package type_config

type ConfigTypeIntegration struct {
	IntegrationName     string `yaml:"integrationname"`
	PlatformID          int    `yaml:"platformid"`
	ClickUpListId       string `yaml:"clickuplistid"`
	ClickUpCustomerList string `yaml:"clickupcustomerlist"`
}

type ConfigTypeClickup struct {
	User int64 `yaml:"user"`
}

type ConfigTypeGeneral struct {
	IntegrationDefault int `yaml:"integrationdefault"`
}

type ConfigType struct {
	ConfigGeneral ConfigTypeGeneral       `yaml:"configGeneral"`
	ConfclickUp   ConfigTypeClickup       `yaml:"configclickUp"`
	Integrations  []ConfigTypeIntegration `yaml:"integrations"`
}