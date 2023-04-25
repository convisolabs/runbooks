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

type ProjectCreateParameters struct {
	Input ProjectCreateInputParameters `json:"input"`
}

type ProjectCreateInputParameters struct {
	CompanyId      int    `json:"companyId"`
	Label          string `json:"label"`
	Goal           string `json:"goal"`
	PlaybooksIds   []int  `json:"playbooksIds"`
	Scope          string `json:"scope"`
	TypeId         int    `json:"typeId"`
	StartDate      string `json:"startDate"`
	EstimatedHours string `json:"estimatedHours"`
}

type ProjectCreateResult struct {
	clientMutationId string
	errors           []string
}

type ListReturn struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ListsReturn struct {
	Lists []ListReturn `json:"lists"`
}

type TaskReturn struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TasksReturn struct {
	Tasks []TaskReturn `json:"tasks"`
}
