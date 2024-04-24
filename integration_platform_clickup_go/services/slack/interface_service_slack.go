package slack_service

import "integration_platform_clickup_go/types/type_slack"

type ISlackService interface {
	RequestPostMessage(request type_slack.PostMessage) error
}
