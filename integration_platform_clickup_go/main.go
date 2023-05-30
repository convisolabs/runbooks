package main

import (
	"bytes"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
	ServiceClickup "integration.platform.clickup/services/service_clickup"
	ServiceConvisoPlatform "integration.platform.clickup/services/service_conviso_platform"
	TypeClickup "integration.platform.clickup/types/type_clickup"
	TypePlatform "integration.platform.clickup/types/type_platform"
	Functions "integration.platform.clickup/utils/functions"
	VariablesConstant "integration.platform.clickup/utils/variables_constant"
	VariablesGlobal "integration.platform.clickup/utils/variables_global"
)

const BANNER = `
____  _       _    __                       ____ _ _      _    _   _       
|  _ \| | __ _| |_ / _| ___  _ __ _ __ ___  / ___| (_) ___| | _| | | |_ __  
| |_) | |/ _∎ | __| |_ / _ \| '__| '_ ∎ _ \| |   | | |/ __| |/ / | | | '_ \ 
|  __/| | (_| | |_|  _| (_) | |  | | | | | | |___| | | (__|   <| |_| | |_) |
|_|   |_|\__,_|\__|_|  \___/|_|  |_| |_| |_|\____|_|_|\___|_|\_\\___/| .__/ 
																	 |_|    
`

func LoadProjects() {
	projects := Functions.LoadCustomerByYamlFile()

	fmt.Println("------Projets------")
	// Print the data
	for i := 0; i < len(projects); i++ {
		fmt.Println(i, " - ", projects[i].IntegrationName)
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
		fmt.Println("Project Selected: ", VariablesGlobal.Customer.IntegrationName)
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

func UpdateClickUpConvisoPlatform(justVerify bool) {
	fmt.Println("...Starting ClickUp Automation...")

	lists, err := ServiceClickup.ReturnLists()

	if err != nil {
		fmt.Println("Error return list: ", err.Error())
		return
	}

	fmt.Println("...Searching valid list...")
	for i := 0; i < len(lists.Lists); i++ {

		fmt.Println("Found List ", lists.Lists[i].Name)

		if Functions.CustomerExistsYamlFileByClickUpListId(lists.Lists[i].Id, Functions.LoadCustomerByYamlFile()) {

			var sliceEpicId []string

			fmt.Println("Found valid list ", lists.Lists[i].Name)

			time.Sleep(time.Second)

			tasks, err := ServiceClickup.ReturnTasks(lists.Lists[i].Id)

			if err != nil {
				fmt.Println("Error return tasks: ", err.Error())
				return
			}

			for j := 0; j < len(tasks.Tasks); j++ {
				fmt.Println("Task ", j+1, "/", len(tasks.Tasks), " - ", tasks.Tasks[j].Name)

				auxEpicTaskId := ""

				if len(tasks.Tasks[j].LinkedTasks) == 0 {
					fmt.Println("Error 0 epics", " :: ", lists.Lists[i].Name, " - ", tasks.Tasks[j].Name)
					continue
				}

				if len(tasks.Tasks[j].LinkedTasks) > 1 {
					fmt.Println("Error 2 epics:", " :: ", lists.Lists[i].Name, " - ", tasks.Tasks[j].Name)
					continue
				}

				//dependendo a ordem que vc linkar as tarefas ele vai jogar no linkid ou no taskid
				if tasks.Tasks[j].Id == tasks.Tasks[j].LinkedTasks[0].TaskId {
					auxEpicTaskId = tasks.Tasks[j].LinkedTasks[0].LinkId
				} else {
					auxEpicTaskId = tasks.Tasks[j].LinkedTasks[0].TaskId
				}

				if slices.Contains(sliceEpicId, auxEpicTaskId) {
					continue
				}

				sliceEpicId = append(sliceEpicId, auxEpicTaskId)

				time.Sleep(time.Second)

				taskEpic, err := ServiceClickup.ReturnTask(auxEpicTaskId)
				if err != nil {
					fmt.Println("Error return task: ", err.Error())
					return
				}

				if justVerify {
					time.Sleep(time.Second)
					ServiceClickup.VerifyTasks(taskEpic)
				} else {
					time.Sleep(time.Second)

					allSubTasksDone := true
					var timeSpent int64
					var requestTask TypeClickup.TaskRequest

					for k := 0; k < len(taskEpic.LinkedTasks); k++ {

						auxTaskId := ""

						if taskEpic.Id == taskEpic.LinkedTasks[k].LinkId {
							auxTaskId = taskEpic.LinkedTasks[k].TaskId
						} else {
							auxTaskId = taskEpic.LinkedTasks[k].LinkId
						}

						taskAux, err := ServiceClickup.ReturnTask(auxTaskId)
						if err != nil {
							fmt.Println("Error taskAux: " + err.Error())
							continue
						}
						auxDuoDate, _ := strconv.ParseInt(taskAux.DueDate, 10, 64)
						auxStartDate, _ := strconv.ParseInt(taskAux.StartDate, 10, 64)
						requestTask.TimeEstimate = requestTask.TimeEstimate + taskAux.TimeEstimate
						timeSpent = timeSpent + taskAux.TimeSpent
						if auxDuoDate > requestTask.DueDate {
							requestTask.DueDate = auxDuoDate
						}

						if auxStartDate != 0 && auxStartDate < requestTask.StartDate || requestTask.StartDate == 0 {
							requestTask.StartDate = auxStartDate
						}

						if taskAux.Status.Status != "done" && taskAux.Status.Status != "canceled" {
							allSubTasksDone = false
						}

						requestTask.Status = ServiceClickup.RetNewStatus(taskEpic.Status.Status, taskAux.Status.Status)
						taskEpic.Status.Status = requestTask.Status

						//precisa alterar o requirements
						//https://app.convisoappsec.com/scopes/553/projects/15296/project_requirements/232561

						ServiceConvisoPlatform.ChangeActivitiesStatus(ServiceClickup.RetCustomFieldUrlConviso(taskAux.CustomFields))

					}

					if allSubTasksDone {
						requestTask.Status = "done"
					}

					if (timeSpent - taskEpic.TimeSpent) > 0 {
						var taskTimeSpentRequest TypeClickup.TaskTimeSpentRequest
						taskTimeSpentRequest.Duration = timeSpent - taskEpic.TimeSpent
						taskTimeSpentRequest.Start = time.Now().UTC().UnixMilli()
						taskTimeSpentRequest.TaskId = taskEpic.Id
						ServiceClickup.RequestTaskTimeSpent(taskEpic.TeamId, taskTimeSpentRequest)
					}

					err := ServiceClickup.RequestPutTask(taskEpic.Id, requestTask)

					if err != nil {
						fmt.Println("Error taskAux: " + err.Error())
					}

					//precisa alterar o status no conviso platform

				}
			}
		}
	}
	fmt.Println("...Finishing ClickUp Automation...")
}

// func MenuClickup() {
// 	var input int
// 	for ok := true; ok; ok = (input != 0) {
// 		fmt.Println("-----Menu Clickup-----")
// 		fmt.Println("0 - Previous Menu")
// 		fmt.Println("1 - Verification Tasks Clickup")
// 		fmt.Println("2 - Update Tasks Clickup")
// 		fmt.Print("Enter the option: ")
// 		n, err := fmt.Scan(&input)
// 		if n < 1 || err != nil {
// 			fmt.Println("Invalid Input")
// 		}
// 		switch input {
// 		case 0:
// 			break
// 		case 1:
// 			UpdateClickUpConvisoPlatform(true)
// 		case 2:
// 			UpdateClickUpConvisoPlatform(false)
// 		default:
// 			fmt.Println("Invalid Input")
// 		}
// 	}
// }

func MenuSearchConvisoPlatform() {
	var input int
	for ok := true; ok; ok = (input != 0) {
		fmt.Println("-----Menu Search Conviso Platform-----")
		fmt.Println("0 - Previous Menu")
		fmt.Println("1 - Requirements")
		fmt.Println("2 - Type Project")
		fmt.Print("Enter the option: ")
		n, err := fmt.Scan(&input)
		if n < 1 || err != nil {
			fmt.Println("Invalid Input")
		}
		switch input {
		case 0:
			break
		case 1:
			ServiceConvisoPlatform.InputSearchRequimentsPlatform()
		case 2:
			ServiceConvisoPlatform.InputSearchProjectTypesPlatform()
		default:
			fmt.Println("Invalid Input")
		}
	}
}

func MainMenu() {
	var input int

	for ok := true; ok; ok = (input != 0) {
		fmt.Println("-----Main Menu-----")
		fmt.Println("Project Selected: ", VariablesGlobal.Customer.IntegrationName)
		fmt.Println("0 - Exit")
		fmt.Println("1 - Menu Setup")
		fmt.Println("2 - Create Project Conviso Platform/ClickUp")
		fmt.Println("3 - Menu Search Conviso Platform")

		fmt.Print("Enter the option: ")
		n, err := fmt.Scan(&input)

		if n < 1 || err != nil {
			fmt.Println("Invalid Input")
			break
		}

		switch input {
		case 0:
			fmt.Println("Finished program!")
		// case 1:
		// 	MenuClickup()
		case 1:
			MenuSetupConfig()
		case 2:
			if VariablesGlobal.Customer.PlatformID == 0 {
				fmt.Println("No Project Selected!")
			} else {
				CreateProject()
			}
		case 3:
			MenuSearchConvisoPlatform()
		default:
			fmt.Println("Invalid Input")
		}
	}
}

func CreateProject() {
	playbookIds := ""
	typeId := 10
	SubTaskReqActivies := "n"

	fmt.Print("Label: ")
	label := Functions.GetTextWithSpace()

	fmt.Print("Goal: ")
	goal := Functions.GetTextWithSpace()

	fmt.Print("Scope: ")
	scope := Functions.GetTextWithSpace()

	fmt.Print("TypeId (Consulting = 10): ")
	n, err := fmt.Scan(&typeId)
	if n < 1 || err != nil {
		fmt.Println("Invalid Input")
		return
	}

	fmt.Print("Playbook (1;2;3): ")
	n, err = fmt.Scan(&playbookIds)
	if n < 1 || err != nil {
		fmt.Println("Invalid Input")
		return
	}

	fmt.Print("Create subtasks with requirements activities? (y or n)")
	n, err = fmt.Scan(&SubTaskReqActivies)
	if n < 1 || err != nil {
		fmt.Println("Invalid Input")
		return
	}

	createConvisoPlatform := TypePlatform.ProjectCreateInputRequest{VariablesGlobal.Customer.PlatformID,
		label, goal, scope, typeId,
		Functions.ConvertStringToArrayInt(playbookIds),
		time.Now().Add(time.Hour * 24).Format("2006-01-02"), "1"}

	err = ServiceConvisoPlatform.AddPlatformProject(createConvisoPlatform)

	if err != nil {
		fmt.Println("Error CreateProject: ", err.Error())
	}

	project, err := ServiceConvisoPlatform.ConfirmProjectCreate(VariablesGlobal.Customer.PlatformID, label)

	if err != nil {
		fmt.Println("Erro CreateProject: Contact the system administrator")
		return
	}

	CustomFieldCustomerPosition, err := ServiceClickup.RetCustomerPosition()
	if err != nil {
		fmt.Println("Error CreateProject: CustomerCustomField error...")
		CustomFieldCustomerPosition = ""
	}

	customFieldCustomer := TypeClickup.CustomFieldRequest{
		"4493a404-3ef7-4d7a-91e4-830ebc666353",
		CustomFieldCustomerPosition,
	}

	customFieldUrlConvisoPlatform := TypeClickup.CustomFieldRequest{
		"8e2863f4-e11f-409c-a373-893bc12200fb",
		"https://app.convisoappsec.com/scopes/" + strconv.Itoa(VariablesGlobal.Customer.PlatformID) + "/projects/" + project.Id,
	}

	customFieldsMainTask := []TypeClickup.CustomFieldRequest{
		customFieldUrlConvisoPlatform,
		TypeClickup.CustomFieldRequest{
			"664816bc-a899-45ec-9801-5a1e5be9c5f6",
			"0",
		},
		customFieldCustomer,
	}

	//create main
	taskMainClickup, err := ServiceClickup.TaskCreateRequest(
		TypeClickup.TaskCreateRequest{
			project.Label,
			project.Scope,
			"backlog",
			true,
			"",
			"",
			customFieldsMainTask})

	if err != nil {
		fmt.Println("Error CreateProject: Problem create ClickUpTask :: ", err.Error())
		return
	}

	if strings.ToLower(SubTaskReqActivies) == "y" {

		for i := 0; i < len(project.Activities); i++ {
			var convisoPlatformUrl bytes.Buffer
			convisoPlatformUrl.WriteString(VariablesConstant.CONVISO_PLATFORM_URL_BASE)
			convisoPlatformUrl.WriteString("scopes/")
			convisoPlatformUrl.WriteString(strconv.Itoa(VariablesGlobal.Customer.PlatformID))
			convisoPlatformUrl.WriteString("/projects/")
			convisoPlatformUrl.WriteString(project.Id)
			convisoPlatformUrl.WriteString("/project_requirements/")
			convisoPlatformUrl.WriteString(project.Activities[i].Id)

			customFieldUrlConvisoPlatformSubTask := TypeClickup.CustomFieldRequest{
				"8e2863f4-e11f-409c-a373-893bc12200fb",
				convisoPlatformUrl.String(),
			}

			customFieldsSubTask := []TypeClickup.CustomFieldRequest{
				customFieldUrlConvisoPlatformSubTask,
				TypeClickup.CustomFieldRequest{
					"664816bc-a899-45ec-9801-5a1e5be9c5f6",
					"2",
				},
				customFieldCustomer}

			_, err := ServiceClickup.TaskCreateRequest(
				TypeClickup.TaskCreateRequest{
					project.Activities[i].Title,
					project.Activities[i].Description,
					"backlog",
					true,
					taskMainClickup.Id,
					taskMainClickup.Id,
					customFieldsSubTask})
			if err != nil {
				fmt.Println("Error CreateProject: Problem create ClickUp SubTask ", project.Activities[i].Title)
			}
		}
	}

	fmt.Println("Create Task Success!")
}

func main() {
	//próximas tarefas
	// qdo atualizar um epico atualizar o projeto no conviso platform
	// qdo atualizar uma task atualizar também o requirements na platforma
	// sanitizar a saída do conviso platform está colocando tags html no clickup

	integrationJustVerify := flag.Bool("iv", false, "Verify if clickup tasks is ok")
	integrationUpdateTasks := flag.Bool("iu", false, "Update Conviso Platform and ClickUp Tasks")
	deploy := flag.Bool("d", false, "See info about deploys")

	flag.Parse()

	if *integrationJustVerify {
		UpdateClickUpConvisoPlatform(true)
	}

	if *integrationUpdateTasks {
		UpdateClickUpConvisoPlatform(false)
	}

	if *deploy {
		ServiceConvisoPlatform.RetDeploys()
	}

	if !*integrationJustVerify && !*integrationUpdateTasks && !*deploy {
		MainMenu()
	}

}
