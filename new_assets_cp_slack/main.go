package main

import (
	"bytes"
	"flag"
	"fmt"
	crawler_service "new_assets_cp_slack/services/crawler"
	type_config "new_assets_cp_slack/types/config"
	"new_assets_cp_slack/utils/constants"
	"new_assets_cp_slack/utils/functions"
	"new_assets_cp_slack/utils/globals"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	fmt.Println("Crawler CP Slack")

	crawlerExecute := flag.Bool("e", false, "Search New Assets Fortify Integration Conviso Platform")
	version := flag.Bool("v", false, "Crawler CP Slack Version")

	flag.Parse()

	if *version {
		fmt.Println("Crawler CP Slack Version: ", constants.VERSION)
		Exit()
	}

	if !InitialCheck() {
		fmt.Println("You need to correct the above information before rerunning the application")
		fmt.Println("Press the Enter Key to finish!")
		fmt.Scanln()
		os.Exit(0)
	}

	SetDefaultValue()

	//iniciando singletons
	crawler_service.GetCrawlerServiceSingletonInstance()
	//fim iniciando singletons

	if *crawlerExecute {
		for i := 0; i < len(globals.Config.Integrations); i++ {
			fmt.Println("Found List ", globals.Config.Integrations[i].IntegrationName)
			AssetsNewCrawlerFortifyIntegration(globals.Config.Integrations[i])
			fmt.Println("Begin: ", time.Now().Format("2006-01-02 15:04:05"))
		}
		os.Exit(0)
	}

	Exit()
}

func SetDefaultValue() {
	if globals.Config.ConfigSlack.HttpAttempt == nil {
		*globals.Config.ConfigSlack.HttpAttempt = 3
	}
}

func Exit() {
	fmt.Println("Finished Crawler CP Slack")
	fmt.Println("Press the Enter Key to finish!")
	fmt.Scanln()
	os.Exit(0)
}

func InitialCheck() bool {
	ret := true

	err := error(nil)

	globals.Config, err = functions.LoadConfigsByYamlFile()

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
		cont := crawlerService.Exec(integration.PlatformID, urlPage)
		if !cont {
			break
		}
		page++
	}
}
