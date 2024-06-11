package cp_service

import "integration_platform_clickup_go/types/type_platform"

type ICPService interface {
	AddPlatformProject(inputParameters type_platform.ProjectCreateRequestInput) (type_platform.ProjectCreateResponse, error)
	SearchRequimentsPlatform(reqSearch string)
	InputSearchRequimentsPlatform()
	InputSearchProjectTypesPlatform()
}
