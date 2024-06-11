package type_platform

type MetadataResponse struct {
	CurrentPage int `json:"currentPage"`
	TotalPages  int `json:"totalPages"`
}

// type JsonData struct {
// 	Data ProjectCreateDataResponse `json:"data"`
// }

type RequirementsParameters struct {
	CompanyId, Page int
	Requirement     string
}

type RequirementsResponse struct {
	Data struct {
		Playbooks struct {
			Collection []Requirements   `json:"collection"`
			Metadata   MetadataResponse `json:"metadata"`
		} `json:"playbooks"`
	} `json:"data"`
}

type Requirements struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

type ProjectCreateRequest struct {
	Input ProjectCreateRequestInput `json:"input"`
}

type ProjectCreateRequestInput struct {
	CompanyId      int    `json:"companyId"`
	Label          string `json:"label"`
	Goal           string `json:"goal"`
	Scope          string `json:"scope"`
	TypeId         int    `json:"typeId"`
	PlaybooksIds   []int  `json:"playbooksIds"`
	StartDate      string `json:"startDate"`
	EstimatedHours string `json:"estimatedHours"`
}

type ProjectCreateResponse struct {
	Data struct {
		CreateProject struct {
			Errors  []string `json:"errors"`
			Project Project  `json:"project"`
		} `json:"createProject"`
	} `json:"data"`
}

// type ProjectCreateReponse struct {
// 	Data ProjectCreateDataResponse `json:"data"`
// }

// type ProjectCreateDataResponse struct {
// 	Playbooks ProjectCreateCreateProjectResponse `json:"createProject"`
// }

// type ProjectCreateCreateProjectResponse struct {
// 	Errors  []string
// 	Project ProjectCollectionResponse
// }

type Activity struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type Company struct {
	Id string `json:"id"`
}

// type Project struct {
// 	Id         string                       `json:"id"`
// 	Label      string                       `json:"label"`
// 	Objective  string                       `json:"objective"`
// 	Scope      string                       `json:"scope"`
// 	Company    CompanyCollectionResponse    `json:"company"`
// 	Activities []ActivityCollectionResponse `json:"activities"`
// }

// type ProjectsResponse struct {
// 	Data ProjectsDataResponse `json:"data"`
// }

// type ProjectsDataResponse struct {
// 	Projects ProjectsProjectsResponse `json:"projects"`
// }

// type ProjectsProjectsResponse struct {
// 	Collection []ProjectCollectionResponse `json:"collection"`
// 	Metadata   MetadataResponse            `json:"metadata"`
// }

type ProjectType struct {
	Id          string `json:"id"`
	Label       string `json:"label"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type ProjectTypesResponse struct {
	Data struct {
		ProjectTypes struct {
			Collection []ProjectType    `json:"collection"`
			Metadata   MetadataResponse `json:"metadata"`
		} `json:"projectTypes"`
	} `json:"data"`
}

// type ProjectTypesDataResponse struct {
// 	ProjectTypes ProjectTypesCollectionMetadataResponse `json:"projectTypes"`
// }

// type ProjectTypesCollectionMetadataResponse struct {
// 	Collection []ProjectType    `json:"collection"`
// 	Metadata   MetadataResponse `json:"metadata"`
// }

type ProjectTypeParameters struct {
	Page        int
	ProjectType string
}

// type DeployType struct {
// 	ChangedLines  int         `json:"changedLines"`
// 	NewLines      int         `json:"newLines"`
// 	RemovedLines  int         `json:"removedLines"`
// 	Reviewed      bool        `json:"reviewed"`
// 	CurrentCommit string      `json:"currentCommit"`
// 	Project       ProjectType `json:"project"`
// }

// type DeployTypeCollectionMetadataResponse struct {
// 	Collection []DeployType     `json:"collection"`
// 	Metadata   MetadataResponse `json:"metadata"`
// }

// type DeployTypeDataResponse struct {
// 	DeployTypeData DeployTypeCollectionMetadataResponse `json:"deploysByCompanyId"`
// }

// type DeployTypeResponse struct {
// 	Data DeployTypeDataResponse `json:"data"`
// }

// type PageParameters struct {
// 	Page int `json:"Page"`
// }

type Project struct {
	Id         string     `json:"id"`
	Label      string     `json:"label"`
	Objective  string     `json:"objective"`
	Scope      string     `json:"scope"`
	Company    Company    `json:"company"`
	Activities []Activity `json:"activities"`
	Status     string     `json:"status"`
}

// type ProjectResponse struct {
// 	Data ProjectDataResponse `json:"data"`
// }

// type ProjectDataResponse struct {
// 	Project Project `json:"project"`
// }

// type UpdateRequirementsActivityStatusRequest struct {
// 	Input UpdateRequirementsActivityStatusInputRequest `json:"input"`
// }

// type UpdateRequirementsActivityStatusInputRequest struct {
// 	ActivityId int `json:"activityId"`
// }
