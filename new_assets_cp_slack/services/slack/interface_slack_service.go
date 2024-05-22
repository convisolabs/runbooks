package slack_service

import type_slack "new_assets_cp_slack/types/slack"

type ISlackService interface {
	RequestPostMessage(request type_slack.PostMessage) error
}
