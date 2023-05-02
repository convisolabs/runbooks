package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	ServicesClickup "integration.platform.clickup/services/services_clickup"
	TypesPlatform "integration.platform.clickup/types/types_platform"
	VariablesGlobal "integration.platform.clickup/utils/variables_global"
)

// criar projeto no platform
// consultor projeto criado e pegar os requirements
// criar projeto no clickup com subtasks do requirements

const CONVISO_PLATFORM_REQUIREMENTS_QUERY = `
	query Playbooks($CompanyId:ID!,$Requirement:String,$Page:Int){
	  playbooks(id: $CompanyId, page: $Page, limit: 10, params: {
	    labelCont:$Requirement
	  }) {
	    collection {
	      checklistTypeId
	      companyId
	      createdAt
	      deletedAt
	      description
	      id
	      label
	      updatedAt
	    }
	    metadata {
	      currentPage
	      limitValue
	      totalCount
	      totalPages
	    }
	  }
	}
`

const CONVISO_PLATFORM_PROJECT_CREATE = `
	mutation CreateProject($input:CreateProjectInput!)
	{
		createProject(
		input: $input
		) 
		{
			clientMutationId
			errors
			project 
			{
				apiCode
				apiResponseReview
				closeComments
				companyId
				connectivity
				continuousDelivery
				contractedHours
				createdAt
				deploySendFrequency
				dueDate
				endDate
				environmentInvaded
				estimatedDays
				estimatedHours
				executiveSummary
				freeRetest
				hasOpenRetest
				hoursOrDays
				id
				integrationDeploy
				inviteToken
				isOpen
				isPublic
				label
				language
				mainRecommendations
				microserviceFolder
				negativeScope
				notificationList
				objective
				pid
				plannedStartedAt
				playbookFinishedAt
				playbookStartedAt
				receiveDeploys
				repositoryUrl
				sacCode
				sacProjectId
				scope
				secretId
				sshPublicKey
				startDate
				status
				students
				subScopeId
				totalAnalysisLines
				totalChangedLines
				totalNewLines
				totalPublishedVulnerabilities
				totalRemovedLines
				type
				updatedAt
				userableId
				userableType
				waiting
			}
		}
	}
`

const BANNER = `
____  _       _    __                       ____ _ _      _    _   _       
|  _ \| | __ _| |_ / _| ___  _ __ _ __ ___  / ___| (_) ___| | _| | | |_ __  
| |_) | |/ _∎ | __| |_ / _ \| '__| '_ ∎ _ \| |   | | |/ __| |/ / | | | '_ \ 
|  __/| | (_| | |_|  _| (_) | |  | | | | | | |___| | | (__|   <| |_| | |_) |
|_|   |_|\__,_|\__|_|  \___/|_|  |_| |_| |_|\____|_|_|\___|_|\_\\___/| .__/ 
																	 |_|    
`

//dog = strings.ReplaceAll(dog, "∎", "`")

func LoadProjects() {
	// Read the file
	data, err := ioutil.ReadFile("projects.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a struct to hold the YAML data
	var projects []VariablesGlobal.CustomerType

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(data, &projects)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("------Projets------")
	// Print the data
	for i := 0; i < len(projects); i++ {
		fmt.Println(i, " - ", projects[i].Name)
	}

	var input int
	fmt.Print("Enter the option: ")
	n, err := fmt.Scan(&input)
	if n < 1 || err != nil || input > len(projects)-1 {
		fmt.Println("Invalid Input")
		return
	}

	VariablesGlobal.Customer = projects[input]

}

func MenuSetupConfig() {
	var input int
	for ok := true; ok; ok = (input != 0) {
		fmt.Println("-----Menu Config-----")
		fmt.Println("Project Selected: ", VariablesGlobal.Customer.Name)
		fmt.Println("0 - Previous Menu")
		fmt.Println("1 - Choose Project Work")
		fmt.Print("Enter the option: ")
		n, err := fmt.Scan(&input)
		if n < 1 || err != nil {
			fmt.Println("Invalid Input")
			break
		}
		switch input {
		case 0:
			break
		case 1:
			LoadProjects()
		default:
			fmt.Println("Invalid Input")
		}
	}
}

func MenuRequirementsSearch() {
	var input int
	for ok := true; ok; ok = (input != 0) {
		fmt.Println("-----Menu Requirements Search-----")
		fmt.Println("Project Selected: ", VariablesGlobal.Customer.Name)
		fmt.Println("0 - Previous Menu")
		fmt.Println("1 - Search Requirements")
		fmt.Print("Enter the option: ")
		n, err := fmt.Scan(&input)
		if n < 1 || err != nil {
			fmt.Println("Invalid Input")
			break
		}
		switch input {
		case 0:
			break
		case 1:
			SearchRequimentsPlatform()
		default:
			fmt.Println("Invalid Input")
		}
	}
}

func SearchRequimentsPlatform() {

	var tokenPlatform = os.Getenv("CONVISO_PLATFORM_TOKEN")
	var reqSearch string
	fmt.Print("Enter part of the requirement: ")
	n, err := fmt.Scan(&reqSearch)
	if n < 1 || err != nil {
		fmt.Println("Invalid Input")
		return
	}

	var result TypesPlatform.RequirementsReturnType
	for i := 0; i <= result.Data.Playbooks.Metadata.TotalPages; i++ {

		parameters, _ := json.Marshal(VariablesGlobal.RequirementsParametersType{CompanyId: VariablesGlobal.Customer.PlatformID, Page: i + 1, Requirement: reqSearch})
		body, _ := json.Marshal(map[string]string{
			"query":     CONVISO_PLATFORM_REQUIREMENTS_QUERY,
			"variables": string(parameters),
		})

		payload := bytes.NewBuffer(body)
		req, err := http.NewRequest(http.MethodPost, "https://app.convisoappsec.com/graphql", payload)
		if err != nil {
			fmt.Println("Error")
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("x-api-key", tokenPlatform)
		client := &http.Client{Timeout: time.Second * 10}
		resp, err := client.Do(req)
		defer req.Body.Close()
		if err != nil {
			fmt.Println("Error")
		}
		data, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal([]byte(string(data)), &result)

		fmt.Println("Results - Requeriments - Project ", VariablesGlobal.Customer.Name)
		for i := 0; i < len(result.Data.Playbooks.Collection)-1; i++ {
			fmt.Println("Requiment ID: ", result.Data.Playbooks.Collection[i].Id, "; Requirement Name:", result.Data.Playbooks.Collection[i].Label)
		}

		if result.Data.Playbooks.Metadata.CurrentPage != result.Data.Playbooks.Metadata.TotalPages {
			var input int
			fmt.Print("See next results? 0 - no; 1 - yes")
			fmt.Print("Enter the option: ")
			fmt.Scan(&input)
			if input == 0 {
				break
			}
		}
		result.Data.Playbooks.Metadata.TotalPages = result.Data.Playbooks.Metadata.TotalPages - 1
	}

}

func MainMenu() {
	var input int

	if slices.Contains(os.Args, "-clickupautomation") {
		ServicesClickup.ClickUpAutomation(false)
	} else {

		for ok := true; ok; ok = (input != 0) {
			fmt.Println("-----Main Menu-----")
			fmt.Println("Project Selected: ", VariablesGlobal.Customer.Name)
			fmt.Println("0 - Exit")
			fmt.Println("1 - Atualizar Projetos ClickUp")
			fmt.Println("2 - Verificar Projetos ClickUp")
			// fmt.Println("1 - Menu Setup")
			// fmt.Println("2 - Menu Search Requirements Conviso Platform")
			// fmt.Println("3 - Create Project Conviso Platform/ClickUp")

			fmt.Print("Enter the option: ")
			n, err := fmt.Scan(&input)

			if n < 1 || err != nil {
				fmt.Println("Invalid Input")
				break
			}

			switch input {
			case 0:
				fmt.Println("Finished program!")
			case 1:
				ServicesClickup.ClickUpAutomation(false)
			case 2:
				ServicesClickup.ClickUpAutomation(true)
			// case 1:
			// 	MenuSetupConfig()
			// case 2:
			// 	MenuRequirementsSearch()
			default:
				fmt.Println("Invalid Input")
			}
		}
	}
}

func AddPlatformProject() {

	var tokenPlatform = os.Getenv("CONVISO_PLATFORM_TOKEN")

	projectParameters := TypesPlatform.ProjectCreateParameters{TypesPlatform.ProjectCreateInputParameters{553, "teste tiago", "teste tiago", []int{475}, "teste", 10, "2023-04-16", "10"}}

	parameters, _ := json.Marshal(projectParameters)
	body, _ := json.Marshal(map[string]string{
		"query":     CONVISO_PLATFORM_PROJECT_CREATE,
		"variables": string(parameters),
	})

	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, "https://app.convisoappsec.com/graphql", payload)
	if err != nil {
		fmt.Println("Error")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", tokenPlatform)
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		fmt.Println("Error")
	}
	data, _ := ioutil.ReadAll(resp.Body)

	var result TypesPlatform.ProjectCreateResult

	json.Unmarshal([]byte(string(data)), &result)

	fmt.Println("Results - Requeriments - Project ", VariablesGlobal.Customer.Name)

}

func main() {
	MainMenu()
	//AddPlatformProject()
}
