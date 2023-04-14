package TypesPlatform

type RequirementsReturnType struct {
	Data RequirementsReturnDataType `json:"data"`
}

type RequirementsReturnDataType struct {
	Playbooks RequirementsReturnPlaybooksType `json:"playbooks"`
}

type RequirementsReturnPlaybooksType struct {
	Collection []RequirementsReturnCollection `json:"collection"`
	Metadata   RequirementsReturnMetadataType `json:"metadata"`
}

type RequirementsReturnMetadataType struct {
	CurrentPage int `json:"currentPage"`
	TotalPages  int `json:"totalPages"`
}

type RequirementsReturnCollection struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}
