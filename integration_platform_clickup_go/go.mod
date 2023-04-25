module main

go 1.20

replace utils/globals => /Users/zani0x03/dev/runbooks/integration_platform_clickup_go/utils

replace integration.platform.clickup/types => /Users/zani0x03/dev/runbooks/integration_platform_clickup_go/types

require (
	golang.org/x/exp v0.0.0-20230420155640-133eef4313cb // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	integration.platform.clickup/types v0.0.0-00010101000000-000000000000 // indirect
	utils/globals v0.0.0-00010101000000-000000000000 // indirect
)
