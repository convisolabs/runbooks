package cp_service

import (
	"integration_platform_clickup_go/types/type_clickup"
	"integration_platform_clickup_go/types/type_platform"
)

type ICPService interface {
	AddPlatformProject(inputParameters type_platform.ProjectCreateRequestInput) (type_platform.ProjectCreateResponse, error)
	SearchRequimentsPlatform(reqSearch string)
	InputSearchRequimentsPlatform()
	InputSearchProjectTypesPlatform()
	RetProjectIdCustomField(text string) (int, error)
	GetProject(id int) (type_platform.Project, error)
	UpdateActivityRequirement(task type_clickup.TaskResponse, project type_platform.Project) error
	UpdateProjectRest(request type_clickup.TaskRequestStore, cpProjectId string, timeEstimate int64) error
}
