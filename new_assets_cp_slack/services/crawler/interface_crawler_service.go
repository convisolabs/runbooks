package crawler_service

type ICrawlerService interface {
	Exec(company int, url string) bool
}
