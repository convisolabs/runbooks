package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"utils/globals"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	TypesPlatform "integration.platform.clickup/types"
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
	var projects []globals.CustomerType

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

	globals.Customer = projects[input]

}

func MenuSetupConfig() {
	var input int
	for ok := true; ok; ok = (input != 0) {
		fmt.Println("-----Menu Config-----")
		fmt.Println("Project Selected: ", globals.Customer.Name)
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
		fmt.Println("Project Selected: ", globals.Customer.Name)
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

		parameters, _ := json.Marshal(globals.RequirementsParametersType{CompanyId: globals.Customer.PlatformID, Page: i + 1, Requirement: reqSearch})
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

		fmt.Println("Results - Requeriments - Project ", globals.Customer.Name)
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

func ReturnTask(taskId string) TypesPlatform.TaskReturn {
	var urlGetTask bytes.Buffer
	urlGetTask.WriteString("https://api.clickup.com/api/v2/task/")
	urlGetTask.WriteString(taskId)

	req, err := http.NewRequest(http.MethodGet, urlGetTask.String(), nil)
	if err != nil {
		// handle error
	}

	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error")
	}
	data, _ := ioutil.ReadAll(resp.Body)

	var task TypesPlatform.TaskReturn
	json.Unmarshal([]byte(string(data)), &task)

	return task

}

func RetNewStatus(statusPrincipal string, statusLinked string) string {

	newReturn := statusPrincipal

	switch statusLinked {
	case "backlog":
		break
	case "to do":
		if statusPrincipal == "backlog" {
			newReturn = "to do"
		}
		break
	case "in progress", "done":
		if statusPrincipal == "backlog" || statusPrincipal == "to do" || statusPrincipal == "blocked" {
			newReturn = "in progress"
		}
	}
	return newReturn
}

func RequestPutTask(taskId string, request TypesPlatform.TaskRequest) TypesPlatform.TaskReturn {

	var urlPutTask bytes.Buffer
	urlPutTask.WriteString("https://api.clickup.com/api/v2/task/")
	urlPutTask.WriteString(taskId)

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPut, urlPutTask.String(), payload)
	if err != nil {
		fmt.Println("Error")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		fmt.Println("Error")
	}
	data, _ := ioutil.ReadAll(resp.Body)

	var result TypesPlatform.TaskReturn

	json.Unmarshal([]byte(string(data)), &result)

	return result
}

func RequestTaskTimeSpent(teamId string, request TypesPlatform.TaskTimeSpentRequest) TypesPlatform.TaskReturn {

	var urlTaskTimeSpent bytes.Buffer
	urlTaskTimeSpent.WriteString("https://api.clickup.com/api/v2/team/")
	urlTaskTimeSpent.WriteString(teamId)
	urlTaskTimeSpent.WriteString("/time_entries")

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, urlTaskTimeSpent.String(), payload)
	if err != nil {
		fmt.Println("Error")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		fmt.Println("Error")
	}
	data, _ := ioutil.ReadAll(resp.Body)

	var result TypesPlatform.TaskReturn

	json.Unmarshal([]byte(string(data)), &result)

	return result
}

func ClickUpAutomation() {
	// listar todos os projetos/clientes
	// consultar todas as tasks alteradas nas últimas 24hs e com o type task = 2 or 1 history
	// com a task vamos chegar na task principal
	//atualizar duedate
	//atualizar duedate
	//depois timetracked

	var result TypesPlatform.ListsReturn

	req, err := http.NewRequest(http.MethodGet, "https://api.clickup.com/api/v2/folder/114948796/list?archived=false", nil)
	if err != nil {
		// handle error
	}

	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error")
	}
	data, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal([]byte(string(data)), &result)

	for i := 0; i < len(result.Lists); i++ {
		if strings.ToLower(result.Lists[i].Name) == "testeprojetosconsulting" {
			var urlGetTasks bytes.Buffer
			urlGetTasks.WriteString("https://api.clickup.com/api/v2/list/")
			urlGetTasks.WriteString(result.Lists[i].Id)
			urlGetTasks.WriteString("/task?custom_fields=[")
			urlGetTasks.WriteString("{\"field_id\":\"664816bc-a899-45ec-9801-5a1e5be9c5f6\",\"operator\":\">=\",\"value\":\"1\"}")
			urlGetTasks.WriteString("]")
			urlGetTasks.WriteString("&include_closed=true")

			fmt.Println(urlGetTasks.String())

			req, err := http.NewRequest(http.MethodGet, urlGetTasks.String(), nil)
			if err != nil {
				// handle error
			}

			req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
			client := &http.Client{Timeout: time.Second * 10}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error")
			}
			data, _ := ioutil.ReadAll(resp.Body)

			var resultTasks TypesPlatform.TasksReturn
			json.Unmarshal([]byte(string(data)), &resultTasks)

			for j := 0; j < len(resultTasks.Tasks); j++ {
				//vou achar a task principal
				taskPrincipal := ReturnTask(resultTasks.Tasks[j].LinkedTasks[0].TaskId)
				var taskAux TypesPlatform.TaskReturn
				allSubTasksDone := true
				var timeSpent int64

				var requestTask TypesPlatform.TaskRequest

				for k := 0; k < len(taskPrincipal.LinkedTasks); k++ {
					taskAux = ReturnTask(taskPrincipal.LinkedTasks[k].LinkId)
					auxDuoDate, _ := strconv.ParseInt(taskAux.DueDate, 10, 64)
					requestTask.TimeEstimate = requestTask.TimeEstimate + taskAux.TimeEstimate
					timeSpent = timeSpent + taskAux.TimeSpent
					if auxDuoDate > requestTask.DueDate {
						requestTask.DueDate = auxDuoDate
					}
					if taskAux.Status.Status != "done" {
						allSubTasksDone = false
					}
				}

				if allSubTasksDone {
					var taskTimeSpentRequest TypesPlatform.TaskTimeSpentRequest
					taskTimeSpentRequest.Duration = timeSpent - taskPrincipal.TimeSpent
					taskTimeSpentRequest.Start = time.Now().UTC().UnixMilli()
					taskTimeSpentRequest.TaskId = taskPrincipal.Id
					requestTask.Status = "done"
					RequestTaskTimeSpent(taskPrincipal.TeamId, taskTimeSpentRequest)
					//atualizar o valor do task
				} else {
					requestTask.Status = RetNewStatus(taskPrincipal.Status.Status, resultTasks.Tasks[j].Status.Status)
				}

				teste := RequestPutTask(taskPrincipal.Id, requestTask)

				println(teste.Id)
				fmt.Println("Sair do For")
			}

			fmt.Println("Sair do sistema")

		}

	}
	//https://api.clickup.com/api/v2/list/900701540171/task?custom_fields=[{"field_id":"664816bc-a899-45ec-9801-5a1e5be9c5f6","operator":"=","value":"2"}, {"field_id":"664816bc-a899-45ec-9801-5a1e5be9c5f6","operator":"=","value":"1"}]
	fmt.Println("ClickUp Automation Started...")
}

func MainMenu() {
	var input int

	fmt.Println(len(os.Args), os.Args)

	if slices.Contains(os.Args, "-clickupautomation") {
		ClickUpAutomation()
	} else {

		for ok := true; ok; ok = (input != 0) {
			fmt.Println("-----Main Menu-----")
			fmt.Println("Project Selected: ", globals.Customer.Name)
			fmt.Println("0 - Exit")
			fmt.Println("1 - Menu Setup")
			fmt.Println("2 - Menu Search Requirements Conviso Platform")
			fmt.Println("3 - Menu Search Requirements Conviso Platform")

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
				MenuSetupConfig()
			case 2:
				MenuRequirementsSearch()
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

	fmt.Println("Results - Requeriments - Project ", globals.Customer.Name)

}

func main() {
	MainMenu()
	//AddPlatformProject()
}
