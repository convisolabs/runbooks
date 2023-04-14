module main

go 1.20

replace utils/globals => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/utils

replace integration.platform.clickup/types => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/types

require (
	gopkg.in/yaml.v3 v3.0.1 // indirect
	integration.platform.clickup/types v0.0.0-00010101000000-000000000000 // indirect
	utils/globals v0.0.0-00010101000000-000000000000 // indirect
)
