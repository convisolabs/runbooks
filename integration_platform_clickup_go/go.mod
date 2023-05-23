module main

go 1.20

replace integration.platform.clickup/types/type_platform => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/types/type_platform

replace integration.platform.clickup/types/type_clickup => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/types/type_clickup

replace integration.platform.clickup/services/service_clickup => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/services/service_clickup

replace integration.platform.clickup/services/service_conviso_platform => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/services/service_conviso_platform

replace integration.platform.clickup/utils/functions => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/utils/functions

replace integration.platform.clickup/utils/variables_global => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/utils/variables_global

replace integration.platform.clickup/utils/variables_constant => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/utils/variables_constant

require (
	golang.org/x/exp v0.0.0-20230420155640-133eef4313cb // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	integration.platform.clickup/services/service_clickup v0.0.0-00010101000000-000000000000 // indirect
	integration.platform.clickup/services/service_conviso_platform v0.0.0-00010101000000-000000000000 // indirect
	integration.platform.clickup/types/type_clickup v0.0.0-00010101000000-000000000000 // indirect
	integration.platform.clickup/types/type_platform v0.0.0-00010101000000-000000000000 // indirect
	integration.platform.clickup/utils/functions v0.0.0-00010101000000-000000000000 // indirect
	integration.platform.clickup/utils/variables_constant v0.0.0-00010101000000-000000000000 // indirect
	integration.platform.clickup/utils/variables_global v0.0.0-00010101000000-000000000000 // indirect
)
