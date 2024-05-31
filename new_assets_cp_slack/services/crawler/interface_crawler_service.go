package crawler_service

type ICrawlerService interface {
	Exec(company int, url string, slackChannel string) bool
}
