package type_platform

type MetadataResponse struct {
	CurrentPage int `json:"currentPage"`
	TotalPages  int `json:"totalPages"`
}

type RequirementsParameters struct {
	CompanyId, Page int
	Requirement     string
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
	Status      string `json:"status"`
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

type ProjectType struct {
	Id          string `json:"id"`
	Label       string `json:"label"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type ProjectTypesResponse struct {
	Data ProjectTypesDataResponse `json:"data"`
}

type ProjectTypesDataResponse struct {
	ProjectTypes ProjectTypesCollectionMetadataResponse `json:"projectTypes"`
}

type ProjectTypesCollectionMetadataResponse struct {
	Collection []ProjectType    `json:"collection"`
	Metadata   MetadataResponse `json:"metadata"`
}

type ProjectTypeParameters struct {
	Page        int
	ProjectType string
}

type DeployType struct {
	ChangedLines  int         `json:"changedLines"`
	NewLines      int         `json:"newLines"`
	RemovedLines  int         `json:"removedLines"`
	Reviewed      bool        `json:"reviewed"`
	CurrentCommit string      `json:"currentCommit"`
	Project       ProjectType `json:"project"`
}

type DeployTypeCollectionMetadataResponse struct {
	Collection []DeployType     `json:"collection"`
	Metadata   MetadataResponse `json:"metadata"`
}

type DeployTypeDataResponse struct {
	DeployTypeData DeployTypeCollectionMetadataResponse `json:"deploysByCompanyId"`
}

type DeployTypeResponse struct {
	Data DeployTypeDataResponse `json:"data"`
}

type PageParameters struct {
	Page int `json:"Page"`
}

type Project struct {
	Id         string                       `json:"id"`
	Label      string                       `json:"label"`
	Objective  string                       `json:"objective"`
	Scope      string                       `json:"scope"`
	Company    CompanyCollectionResponse    `json:"company"`
	Activities []ActivityCollectionResponse `json:"activities"`
	Status     string                       `json:"status"`
}

type ProjectResponse struct {
	Data ProjectDataResponse `json:"data"`
}

type ProjectDataResponse struct {
	Project Project `json:"project"`
}

type UpdateRequirementsActivityStatusRequest struct {
	Input UpdateRequirementsActivityStatusInputRequest `json:"input"`
}

type UpdateRequirementsActivityStatusInputRequest struct {
	ActivityId int `json:"activityId"`
}
