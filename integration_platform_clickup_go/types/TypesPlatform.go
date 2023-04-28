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

type TaskLinkedTasks struct {
	TaskId string `json:"task_id"`
	LinkId string `json:"link_id"`
}

type TaskStatusReturn struct {
	Status string `json:"status"`
}

type TaskReturn struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	Status       TaskStatusReturn  `json:"status"`
	DueDate      string            `json:"due_date"`
	TimeEstimate int64             `json:"time_estimate"`
	SubTasks     []TaskReturn      `json:"subtasks"`
	LinkedTasks  []TaskLinkedTasks `json:"linked_tasks"`
	TeamId       string            `json:"team_id"`
	TimeSpent    int64             `json:"time_spent"`
}

type TasksReturn struct {
	Tasks []TaskReturn `json:"tasks"`
}

type TaskRequest struct {
	DueDate      int64  `json:"due_date"`
	DueDateTime  bool   `json:"due_date_time"`
	TimeEstimate int64  `json:"time_estimate"`
	Status       string `json:"status"`
}

type TaskTimeSpentRequest struct {
	Start    int64  `json:"start"`
	Duration int64  `json:"duration"`
	TaskId   string `json:"tid"`
}
