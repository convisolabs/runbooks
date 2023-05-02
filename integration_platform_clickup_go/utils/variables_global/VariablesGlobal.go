package VariablesGlobal

type CustomerType struct {
	PlatformID    int    `yaml:"platformid"`
	ClickUpListId string `yaml:"clickuplistid"`
	Name          string `yaml:"name"`
}

type RequirementsParametersType struct {
	CompanyId, Page int
	Requirement     string
}

var Customer = CustomerType{0, "", "No selected project"}
