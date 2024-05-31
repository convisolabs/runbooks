package type_config

type ConfigTypeSlack struct {
	HttpAttempt *int `yaml:"httpattempt",omitempty`
}

type ConfigTypeCP struct {
	HttpAttempt *int `yaml:"httpattempt",omitempty`
}

type ConfigTypeIntegration struct {
	IntegrationName  string `yaml:"integrationname"`
	PlatformID       int    `yaml:"platformid"`
	IntegrationType  int    `yaml:"integrationType"`
	SlackChannel     string `yaml:"slackChannel"`
	LastVerification string `yaml:"lastVerification"`
}

type ConfigType struct {
	ConfigSlack  ConfigTypeSlack         `yaml:"configslack"`
	ConfigCP     ConfigTypeCP            `yaml:"configcp"`
	Integrations []ConfigTypeIntegration `yaml:"integrations"`
}
