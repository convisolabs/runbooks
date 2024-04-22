package service_crawler

import (
	"fmt"
	"integration_platform_clickup_go/repositories/repository_assets"
	"integration_platform_clickup_go/services/service_slack"
	"integration_platform_clickup_go/types/type_repository"
	"integration_platform_clickup_go/types/type_slack"
	"integration_platform_clickup_go/utils/variables_constant"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

func Exec(company int, url string) bool {

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

				assetExist, err := repository_assets.AssetExist(asset)

				if err != nil {
					fmt.Println("Service Crawler :: Error :: It was impossible to veriry the asset in database!")
				}

				if !assetExist {
					asset.Id = uuid.New().String()
					asset.DtCreated = time.Now().Format("2006-01-02 15:04:05.000000")
					err = repository_assets.Insert(asset)

					if err != nil {
						fmt.Println("Service Crawler :: Error :: It was impossible to save the asset in database!")
					}

					slackMessage := "*Salve! Salve! <!subteam^S0599SA5YB0|Security Champions da Telefônica Vivo> chegou um Ativo New!*\n" +
						"*Id:* " + asset.CPAssetId + "\n" +
						"*Project:* " + asset.Name + "\n" +
						"*Link:* " + url + "\n\n" +
						"Por favor verifiquem se esse ativo pertence a sua squad!"

					err = service_slack.RequestPostMessage(
						type_slack.PostMessage{
							Channel: variables_constant.SLACK_CHANNEL_CONSULTING,
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
