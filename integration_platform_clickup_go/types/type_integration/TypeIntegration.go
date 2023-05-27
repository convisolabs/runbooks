package TypeIntegration

type CustomerType struct {
	IntegrationName     string `yaml:"integrationname"`
	PlatformID          int    `yaml:"platformid"`
	ClickUpListId       string `yaml:"clickuplistid"`
	ClickUpCustomerList string `yaml:"clickupcustomerlist"`
}
