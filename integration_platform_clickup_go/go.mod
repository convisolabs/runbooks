module main20

go 1.20

replace integration.platform.clickup/types/type_platform => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/types/type_platform

replace integration.platform.clickup/types/type_clickup => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/types/type_clickup

replace integration.platform.clickup/types/type_integration => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/types/type_integration

replace integration.platform.clickup/services/service_clickup => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/services/service_clickup

replace integration.platform.clickup/services/service_conviso_platform => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/services/service_conviso_platform

replace integration.platform.clickup/utils/functions => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/utils/functions

replace integration.platform.clickup/utils/variables_global => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/utils/variables_global

replace integration.platform.clickup/utils/variables_constant => /home/zani0x03/dev/runbooks/integration_platform_clickup_go/utils/variables_constant

require (
	github.com/PuerkitoBio/goquery v1.8.1 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/jaytaylor/html2text v0.0.0-20230321000545-74c2419ad056 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	golang.org/x/exp v0.0.0-20230420155640-133eef4313cb // indirect
	golang.org/x/net v0.9.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	integration.platform.clickup/services/service_clickup v0.0.0-00010101000000-000000000000 // indirect
	integration.platform.clickup/services/service_conviso_platform v0.0.0-00010101000000-000000000000 // indirect
	integration.platform.clickup/types/type_clickup v0.0.0-00010101000000-000000000000 // indirect
	integration.platform.clickup/types/type_integration v0.0.0-00010101000000-000000000000 // indirect
	integration.platform.clickup/types/type_platform v0.0.0-00010101000000-000000000000 // indirect
	integration.platform.clickup/utils/functions v0.0.0-00010101000000-000000000000 // indirect
	integration.platform.clickup/utils/variables_constant v0.0.0-00010101000000-000000000000 // indirect
	integration.platform.clickup/utils/variables_global v0.0.0-00010101000000-000000000000 // indirect
)
