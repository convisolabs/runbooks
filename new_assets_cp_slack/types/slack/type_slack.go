package type_slack

type PostMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}
