package TypesClickup

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

type TaskTimeSpentRequest struct {
	Start    int64  `json:"start"`
	Duration int64  `json:"duration"`
	TaskId   string `json:"tid"`
}