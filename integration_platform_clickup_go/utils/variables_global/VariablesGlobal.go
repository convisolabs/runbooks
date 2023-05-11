package VariablesGlobal

type CustomerType struct {
	PlatformID    int    `yaml:"platformid"`
	ClickUpListId string `yaml:"clickuplistid"`
	Name          string `yaml:"name"`
}

var Customer = CustomerType{0, "", "No selected project"}
