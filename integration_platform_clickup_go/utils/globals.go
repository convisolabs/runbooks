package globals

type CustomerType struct {
	PlatformID    int    `yaml:"platformid"`
	ClickUpListId int    `yaml:"clickuplistid"`
	Name          string `yaml:"name"`
}

type RequirementsParametersType struct {
	CompanyId, Page int
	Requirement     string
}

var Customer = CustomerType{0, 0, "No selected project"}
