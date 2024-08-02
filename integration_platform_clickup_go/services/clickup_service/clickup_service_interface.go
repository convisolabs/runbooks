package clickup_service

import "integration_platform_clickup_go/types/type_clickup"

type IClickupService interface {
	RetClickUpDropDownPosition(clickupListId string, clickupFieldId string, searchValue string) (int, error)
	RetAllCustomFieldByList(listId string) (type_clickup.CustomFieldsResponse, error)
	TaskCreateRequest(request type_clickup.TaskCreateRequest) (type_clickup.TaskResponse, error)
	ReturnList(listId string) (type_clickup.ListResponse, error)
	VerifyErrorsProjectWithStore(list type_clickup.ListResponse)
	UpdateTasksInDoneToClosed(list type_clickup.ListResponse)
	UpdateProjectWithStore(list type_clickup.ListResponse)
	UpdateTasksInDoneToClosedPSHierarchy(list type_clickup.ListResponse, psHierarchy int)
}
