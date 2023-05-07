package TypePlatform

type MetadataResponse struct {
	CurrentPage int `json:"currentPage"`
	TotalPages  int `json:"totalPages"`
}

type RequirementsResponse struct {
	Data RequirementsDataResponse `json:"data"`
}

type RequirementsDataResponse struct {
	Playbooks RequirementsPlaybooksResponse `json:"playbooks"`
}

type RequirementsPlaybooksResponse struct {
	Collection []RequirementsCollectionResponse `json:"collection"`
	Metadata   MetadataResponse                 `json:"metadata"`
}

type RequirementsCollectionResponse struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

type ProjectCreateRequest struct {
	Input ProjectCreateInputRequest `json:"input"`
}

type ProjectCreateInputRequest struct {
	CompanyId      int    `json:"companyId"`
	Label          string `json:"label"`
	Goal           string `json:"goal"`
	Scope          string `json:"scope"`
	TypeId         int    `json:"typeId"`
	PlaybooksIds   []int  `json:"playbooksIds"`
	StartDate      string `json:"startDate"`
	EstimatedHours string `json:"estimatedHours"`
}

type ActivityCollectionResponse struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CompanyCollectionResponse struct {
	Id string `json:"id"`
}

type ProjectCollectionResponse struct {
	Id         string                       `json:"id"`
	Label      string                       `json:"label"`
	Objective  string                       `json:"objective"`
	Scope      string                       `json:"scope"`
	Company    CompanyCollectionResponse    `json:"company"`
	Activities []ActivityCollectionResponse `json:"activities"`
}

type ProjectsResponse struct {
	Data ProjectsDataResponse `json:"data"`
}

type ProjectsDataResponse struct {
	Projects ProjectsProjectsResponse `json:"projects"`
}

type ProjectsProjectsResponse struct {
	Collection []ProjectCollectionResponse `json:"collection"`
	Metadata   MetadataResponse            `json:"metadata"`
}
