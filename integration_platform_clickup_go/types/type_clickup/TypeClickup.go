package TypeClickup

type ListResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ListsResponse struct {
	Lists []ListResponse `json:"lists"`
}

type TaskLinkedTasks struct {
	TaskId string `json:"task_id"`
	LinkId string `json:"link_id"`
}

type TaskStatusResponse struct {
	Status string `json:"status"`
}

type TaskResponse struct {
	Id           string             `json:"id"`
	Name         string             `json:"name"`
	Status       TaskStatusResponse `json:"status"`
	DueDate      string             `json:"due_date"`
	StartDate    string             `json:"start_date"`
	TimeEstimate int64              `json:"time_estimate"`
	SubTasks     []TaskResponse     `json:"subtasks"`
	LinkedTasks  []TaskLinkedTasks  `json:"linked_tasks"`
	TeamId       string             `json:"team_id"`
	TimeSpent    int64              `json:"time_spent"`
	List         ListResponse       `json:"list"`
	CustomFields []CustomField      `json:"custom_fields"`
	Url          string             `json:"url"`
	Parent       string             `json:"parent"`
	Assignees    []AssigneeField    `json:"assignees"`
}

type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

type TaskRequest struct {
	StartDate    int64  `json:"start_date"`
	DueDate      int64  `json:"due_date"`
	TimeEstimate int64  `json:"time_estimate"`
	Status       string `json:"status"`
}

type TaskRequestStore struct {
	StartDate int64  `json:"start_date"`
	DueDate   int64  `json:"due_date"`
	Status    string `json:"status"`
}

type TaskTimeSpentRequest struct {
	Start    int64  `json:"start"`
	Duration int64  `json:"duration"`
	TaskId   string `json:"tid"`
}

type TaskCreateRequest struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Status      string               `json:"status"`
	NotifyAll   bool                 `json:"notify_all"`
	Parent      string               `json:"parent,omitempty"`
	LinksTo     string               `json:"links_to,omitempty"`
	CustomField []CustomFieldRequest `json:"custom_fields"`
}

type CustomFieldRequest struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type CustomFieldsResponse struct {
	Fields []CustomField `json:"fields"`
}

type CustomField struct {
	Id         string                `json:"id"`
	Name       string                `json:"name"`
	TypeConfig CustomFieldTypeConfig `json:"type_config"`
	//ValueString string                `json:"value"`
	Value interface{} `json:"value"`
}

type AssigneeField struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

type CustomFieldTypeConfig struct {
	Options []CustomFieldOptions `json:"options"`
}

type CustomFieldOptions struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	OrderIndex int    `json:"orderindex"`
}
