package type_config

type ConfigTypeSlack struct {
	HttpAttempt *int `yaml:"httpattempt",omitempty`
}

type ConfigTypeIntegration struct {
	IntegrationName string `yaml:"integrationname"`
	PlatformID      int    `yaml:"platformid"`
}

type ConfigType struct {
	ConfigSlack  ConfigTypeSlack         `yaml:"configslack"`
	Integrations []ConfigTypeIntegration `yaml:"integrations"`
}
