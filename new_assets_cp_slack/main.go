package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	cp_service "new_assets_cp_slack/services/cp"
	crawler_service "new_assets_cp_slack/services/crawler"
	slack_service "new_assets_cp_slack/services/slack"
	type_config "new_assets_cp_slack/types/config"
	type_cp "new_assets_cp_slack/types/cp"
	enum_integration_type "new_assets_cp_slack/types/enum/integration_type"
	type_integration "new_assets_cp_slack/types/integration"
	type_slack "new_assets_cp_slack/types/slack"
	"new_assets_cp_slack/utils/constants"
	"new_assets_cp_slack/utils/functions"
	"new_assets_cp_slack/utils/globals"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var iFunc functions.IFunctions
var iCrawlerService crawler_service.ICrawlerService
var iCPService cp_service.ICPService
var iSlackService slack_service.ISlackService

func InitializeDependencyInjection() {
	iFunc = functions.GetFunctionsSingletonInstance()
	iCrawlerService = crawler_service.GetCrawlerServiceSingletonInstance()
	iCPService = cp_service.GetCPServiceSingletonInstance()
	iSlackService = slack_service.GetSlackServiceSingletonInstance()
}

func main() {

	fmt.Println("Integration New Assets")

	execute := flag.Bool("e", false, "Execute the integration")
	version := flag.Bool("v", false, "Integration Version")

	flag.Parse()

	if *version {
		fmt.Println("Integration New Assets Version: ", constants.VERSION)
		Exit()
	}

	InitializeDependencyInjection()

	if !InitialCheck() {
		Exit()
	}

	SetDefaultValue()

	if *execute {

		endDate := time.Now().Format(time.RFC3339)

		for i := 0; i < len(globals.Config.Integrations); i++ {
			fmt.Println("Found List ", globals.Config.Integrations[i].IntegrationName)
			fmt.Println("Begin: ", time.Now().Format("2006-01-02 15:04:05"))
			switch globals.Config.Integrations[i].IntegrationType {
			case enum_integration_type.CRAWLER_FORTIFY:
				AssetsNewCrawlerFortifyIntegration(globals.Config.Integrations[i])
			case enum_integration_type.ASSETS_CP:
				AssetsNewByTimeCP(globals.Config.Integrations[i], endDate)
			}
			globals.Config.Integrations[i].LastVerification = endDate
			fmt.Println("End: ", time.Now().Format("2006-01-02 15:04:05"))
		}

		yamlData, err := yaml.Marshal(&globals.Config)

		if err != nil {
			fmt.Println("main :: Error while Marshaling :: ", err)
			Exit()
		}

		err = iFunc.SaveYamlFile(
			type_integration.SaveFile{
				FileName:    "projects.yaml",
				FileContent: yamlData,
				Perm:        fs.ModePerm,
			},
		)

		if err != nil {
			fmt.Println("main :: it was impossible to save project.yaml :: ", err.Error())
			Exit()
		}
	}

	Exit()
}

func SetDefaultValue() {
	if globals.Config.ConfigSlack.HttpAttempt == nil {
		*globals.Config.ConfigSlack.HttpAttempt = 3
	}

	if globals.Config.ConfigCP.HttpAttempt == nil {
		*globals.Config.ConfigCP.HttpAttempt = 3
	}

	for i := 0; i < len(globals.Config.Integrations); i++ {
		if len(globals.Config.Integrations[i].LastVerification) == 0 {
			globals.Config.Integrations[i].LastVerification = time.Now().Format(time.RFC3339)
		}
	}
}

func Exit() {
	fmt.Println("Finished Integration")
	fmt.Println("Press the Enter Key to continue!")
	fmt.Scanln()
	os.Exit(0)
}

func InitialCheck() bool {
	ret := true

	err := error(nil)

	globals.Config, err = iFunc.LoadConfigsByYamlFile()

	if err != nil {
		fmt.Println("YAML File with Problem! ", err.Error())
		ret = false
	}

	if os.Getenv(constants.SLACK_ASSET_TOKEN_NAME) == "" {
		fmt.Println("Variable ", constants.SLACK_ASSET_TOKEN_NAME, " is empty!")
		ret = false
	}

	if os.Getenv(constants.CONVISO_PLATFORM_TOKEN_NAME) == "" {
		fmt.Println("Variable ", constants.CONVISO_PLATFORM_TOKEN_NAME, " is empty!")
		ret = false
	}

	return ret
}

func AssetsNewCrawlerFortifyIntegration(integration type_config.ConfigTypeIntegration) {

	var urlBase bytes.Buffer
	urlBase.WriteString(constants.CONVISO_PLATFORM_URL_BASE)
	urlBase.WriteString("scopes/")
	urlBase.WriteString(strconv.Itoa(integration.PlatformID))
	urlBase.WriteString("/integrations/fortify/select_projects?page={1}")

	page := 1

	crawlerService := crawler_service.GetCrawlerServiceSingletonInstance()

	for {
		urlPage := strings.Replace(urlBase.String(), "{1}", strconv.Itoa(page), -1)
		cont := crawlerService.Exec(integration.PlatformID, urlPage, integration.SlackChannel)
		if !cont {
			break
		}
		page++
	}
}

func AssetsNewByTimeCP(integration type_config.ConfigTypeIntegration, endDate string) {
	assets, err := iCPService.GetAssetsByTime(type_cp.AssetsByTimeParameters{
		CompanyId: integration.PlatformID,
		Limit:     100,
		Search: type_cp.AssetsByTimeSearchParameters{
			CreatedAt: type_cp.AssetsByTimeCreateAtParameters{
				StartDate: integration.LastVerification,
				EndDate:   endDate,
			},
		},
	})

	if err != nil {
		fmt.Println("Error AssetsNewByTimeCP ", err.Error())
		Exit()
	}

	for i := 0; i < len(assets); i++ {
		slackMessage := "*Salve! Salve! <!subteam^S075C5WR6AH|Security Champions da Universidade Cruzeiro do Sul> chegou um Ativo New!*\n" +
			"*Id:* " + assets[i].Id + "\n" +
			"*Asset:* " + assets[i].Name + "\n\n" +
			"Por favor, procure a Squad e ajudem-os a classificar o ativo e a preencher as informações relevantes na tela de Assets"

		err = iSlackService.RequestPostMessage(
			type_slack.PostMessage{
				Channel: integration.SlackChannel,
				Text:    slackMessage,
			},
		)

		if err != nil {
			fmt.Println("AssetsNewByTimeCP :: it was impossible to send slack message :: ", err.Error())
		}
	}
}
