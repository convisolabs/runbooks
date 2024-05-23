package crawler_service

import (
	"fmt"
	asset_repository "new_assets_cp_slack/repositories/assets"
	slack_service "new_assets_cp_slack/services/slack"
	type_repository "new_assets_cp_slack/types/repository"
	type_slack "new_assets_cp_slack/types/slack"
	"new_assets_cp_slack/utils/constants"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

type CrawlerService struct {
	slackService    slack_service.ISlackService
	assetRepository asset_repository.IAssetReposiroty
}

func CrawlerServiceNew(slackService slack_service.ISlackService, assetRepository asset_repository.IAssetReposiroty) ICrawlerService {
	return &CrawlerService{slackService: slackService, assetRepository: assetRepository}
}

func (f *CrawlerService) Exec(company int, url string) bool {

	ret := true
	c := colly.NewCollector()

	c.OnHTML("div[class]", func(e *colly.HTMLElement) {
		class := e.Attr("class")

		if strings.EqualFold(class, "col-lg-6 col-md-6 col-xs-12") {

			if strings.EqualFold(e.DOM.Find("p").Text(), "new") {

				rProject, _ := regexp.Compile("(.*)Project: ([a-zA-Z0-9\\-\\_ ]*)")
				rId, _ := regexp.Compile("(.*)Id: ([a-zA-Z0-9\\-\\_ ]*)")

				asset := type_repository.Asset{
					CPCompanyId: company,
					Name:        strings.TrimSpace((strings.Split(rProject.FindString(e.DOM.Find("strong").Text()), ":"))[1]),
					CPAssetId:   strings.TrimSpace((strings.Split(rId.FindString(e.DOM.Find("strong").Text()), ":"))[1]),
				}

				assetExist, err := f.assetRepository.AssetExist(asset)

				if err != nil {
					fmt.Println("Service Crawler :: Error :: It was impossible to veriry the asset in database! :: ", err.Error())
				}

				if !assetExist {
					asset.Id = uuid.New().String()
					asset.DtCreated = time.Now().Format("2006-01-02 15:04:05.000000")
					err = f.assetRepository.Insert(asset)

					if err != nil {
						fmt.Println("Service Crawler :: Error :: It was impossible to save the asset in database!")
					}

					slackMessage := "*Salve! Salve! <!subteam^S0599SA5YB0|Security Champions da Telefônica Vivo> chegou um Ativo New!*\n" +
						"*Id:* " + asset.CPAssetId + "\n" +
						"*Project:* " + asset.Name + "\n" +
						"*Link:* " + url + "\n\n" +
						"Por favor verifiquem se esse ativo pertence a sua squad!"

					err = f.slackService.RequestPostMessage(
						type_slack.PostMessage{
							Channel: constants.SLACK_CHANNEL_CONSULTING,
							Text:    slackMessage,
						},
					)

					if err != nil {
						fmt.Println("service_crawler :: it was impossible to send slack message")
					}
				}

			}
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Armature-Api-Key", os.Getenv("CONVISO_PLATFORM_TOKEN"))
		fmt.Println("OnRequest", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.Request, "\nError:", err)
		ret = false
	})

	c.Visit(url)

	return ret
}
