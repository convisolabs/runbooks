package main

import (
	"bytes"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jaytaylor/html2text"
	"golang.org/x/exp/slices"

	ServiceClickup "integration.platform.clickup/services/service_clickup"
	ServiceConvisoPlatform "integration.platform.clickup/services/service_conviso_platform"
	TypeClickup "integration.platform.clickup/types/type_clickup"
	TypeEnumClickupTypeConsulting "integration.platform.clickup/types/type_enum/type_enum_clickup_type_consulting"
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

func VerifyErrorsProjectWithStore(list TypeClickup.ListResponse) {
	VerifySubtask(list, TypeEnumClickupTypeConsulting.EPIC, TypeEnumClickupTypeConsulting.STORE)
	VerifySubtask(list, TypeEnumClickupTypeConsulting.STORE, TypeEnumClickupTypeConsulting.TASK)
	VerifyTasks(list)
}

func VerifyTasks(list TypeClickup.ListResponse) {
	tasks, err := ServiceClickup.ReturnTasks(list.Id, TypeEnumClickupTypeConsulting.TASK)

	if err != nil {
		fmt.Println("Error VerifyTasks :: ", err.Error())
		return
	}

	for i := 0; i < len(tasks.Tasks); i++ {
		task, err := ServiceClickup.ReturnTask(tasks.Tasks[i].Id)

		if err != nil {
			fmt.Println("Error VerifyTasks GetTask :: ", err.Error())
			return
		}

		if task.Parent == "" {
			fmt.Println("TASK Without Store", " :: ", list.Name, " :: ", tasks.Tasks[i].Name, " :: ",
				strings.ToLower(tasks.Tasks[i].Status.Status), " :: ", tasks.Tasks[i].Url,
				" :: ", ServiceClickup.RetAssigness(tasks.Tasks[i].Assignees))
			continue
		}

		if strings.ToLower(task.Status.Status) != "backlog" && strings.ToLower(task.Status.Status) != "canceled" && strings.ToLower(task.Status.Status) != "blocked" {
			if task.DueDate == "" {
				fmt.Println("Task with errors: ", task.List.Name, " - ", task.Name, " - ", task.Name, " :: ", "DueDate empty", " :: ", task.Url,
					" :: ", ServiceClickup.RetAssigness(task.Assignees))
			}

			if task.StartDate == "" {
				fmt.Println("Task with errors: ", task.List.Name, " - ", task.Name, " - ", task.Name, " :: ", "StartDate empty", " :: ", task.Url,
					" :: ", ServiceClickup.RetAssigness(task.Assignees))
			}

			if task.TimeEstimate == 0 {
				fmt.Println("Task with errors: ", task.List.Name, " - ", task.Name, " - ", task.Name, " :: ", "TimeEstimate empty", " :: ", task.Url,
					" :: ", ServiceClickup.RetAssigness(task.Assignees))
			}

			if task.Status.Status == "done" && task.TimeSpent == 0 {
				fmt.Println("Task with errors: ", task.List.Name, " - ", task.Name, " - ", task.Name, " :: ", "TimeSpent empty", " :: ", task.Url,
					" :: ", ServiceClickup.RetAssigness(task.Assignees))
			}
		}
	}
}

func VerifySubtask(list TypeClickup.ListResponse, customFieldTypeConsulting int, customFieldTypeConsultingSubTask int) {

	tasks, err := ServiceClickup.ReturnTasks(list.Id, customFieldTypeConsulting)

	if err != nil {
		fmt.Println("Error VerifySubtask :: ", err.Error())
		return
	}

	for i := 0; i < len(tasks.Tasks); i++ {
		task, err := ServiceClickup.ReturnTask(tasks.Tasks[i].Id)

		if err != nil {
			fmt.Println("Error VerifySubtask GetTask :: ", err.Error())
			return
		}

		if len(task.SubTasks) == 0 {
			fmt.Println(TypeEnumClickupTypeConsulting.ToString(customFieldTypeConsulting),
				" Without ",
				TypeEnumClickupTypeConsulting.ToString(customFieldTypeConsultingSubTask),
				" :: ", list.Name, " :: ", tasks.Tasks[i].Name, " :: ", strings.ToLower(tasks.Tasks[i].Status.Status), " :: ", tasks.Tasks[i].Url,
				" :: ", ServiceClickup.RetAssigness(tasks.Tasks[i].Assignees))
			continue
		}

		for j := 0; j < len(task.SubTasks); j++ {
			subTask, err := ServiceClickup.ReturnTask(task.SubTasks[j].Id)
			if err != nil {
				fmt.Println("Error VerifySubtask GetTask GetSubTask :: ", err.Error())
				return
			}

			customFieldsSubTask := ServiceClickup.RetCustomFieldTypeConsulting(subTask.CustomFields)

			if customFieldsSubTask != customFieldTypeConsultingSubTask {
				fmt.Println(
					subTask.Name,
					" should be ",
					TypeEnumClickupTypeConsulting.ToString(customFieldTypeConsultingSubTask),
					" but is ",
					TypeEnumClickupTypeConsulting.ToString(customFieldsSubTask),
					" :: ", list.Name, " :: ", strings.ToLower(subTask.Status.Status), " :: ", subTask.Url,
					" :: ", ServiceClickup.RetAssigness(subTask.Assignees))
				continue
			}
		}
	}
}

func UpdateProjectWithStore(list TypeClickup.ListResponse) {
	UpdateSubtask(list, TypeEnumClickupTypeConsulting.TASK, TypeEnumClickupTypeConsulting.STORE)
	UpdateSubtask(list, TypeEnumClickupTypeConsulting.STORE, TypeEnumClickupTypeConsulting.EPIC)
}

func UpdateSubtask(list TypeClickup.ListResponse, typeConsultingTask int, typeConsultingParent int) {

	tasks, err := ServiceClickup.ReturnTasks(list.Id, typeConsultingTask)

	if err != nil {
		fmt.Println("Error UpdateSubtask :: ", err.Error())
		return
	}

	var sliceParentId []string

	for i := 0; i < len(tasks.Tasks); i++ {
		if tasks.Tasks[i].Parent == "" {
			continue
		}

		if slices.Contains(sliceParentId, tasks.Tasks[i].Parent) {
			continue
		}

		sliceParentId = append(sliceParentId, tasks.Tasks[i].Parent)

		taskParent, err := ServiceClickup.ReturnTask((tasks.Tasks[i].Parent))

		if err != nil {
			fmt.Println("Error UpdateSubtask GetTask Parent :: ", err.Error())
			continue
		}

		var requestTask TypeClickup.TaskRequestStore
		requestTask.Status = taskParent.Status.Status
		requestTask.DueDate, _ = strconv.ParseInt(taskParent.DueDate, 10, 64)
		requestTask.StartDate, _ = strconv.ParseInt(taskParent.StartDate, 10, 64)
		allTaskDone := true
		hasUpdate := false
		for j := 0; j < len(taskParent.SubTasks); j++ {
			subTask, err := ServiceClickup.ReturnTask(taskParent.SubTasks[j].Id)
			if err != nil {
				fmt.Println("Error UpdateSubtask GetTask SubTask :: ", err.Error())
				return
			}
			var auxStartDate int64
			var auxDueDate int64

			auxStartDate, _ = strconv.ParseInt(subTask.StartDate, 10, 64)
			auxDueDate, _ = strconv.ParseInt(subTask.DueDate, 10, 64)
			if (auxStartDate < requestTask.StartDate) || (auxStartDate != 0 && requestTask.StartDate == 0) {
				requestTask.StartDate = auxStartDate
				hasUpdate = true
			}

			if (auxDueDate > requestTask.DueDate) || (auxDueDate != 0 && requestTask.DueDate == 0) {
				requestTask.DueDate = auxDueDate
				hasUpdate = true
			}

			hasUpdateStatus := false
			requestTask.Status, hasUpdateStatus = ServiceClickup.RetNewStatus(requestTask.Status, subTask.Status.Status)

			if hasUpdateStatus {
				hasUpdate = true
			}

			if subTask.Status.Status != "done" && subTask.Status.Status != "canceled" {
				allTaskDone = false
			}
		}

		if allTaskDone {
			requestTask.Status = "done"
			hasUpdate = true
		}

		if hasUpdate {
			ServiceClickup.RequestPutTaskStore(taskParent.Id, requestTask)
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

	lstCustomersYamlFile := Functions.LoadCustomerByYamlFile()

	for i := 0; i < len(lists.Lists); i++ {

		fmt.Println("Found List ", lists.Lists[i].Name)

		if Functions.CustomerExistsYamlFileByClickUpListId(lists.Lists[i].Id, lstCustomersYamlFile) {

			if justVerify && VariablesGlobal.Customer.HasStore {
				VerifyErrorsProjectWithStore(lists.Lists[i])
				return
			}

			if !justVerify && VariablesGlobal.Customer.HasStore {
				UpdateProjectWithStore(lists.Lists[i])
				return
			}

			var sliceEpicId []string

			fmt.Println("Found valid list ", lists.Lists[i].Name)

			time.Sleep(time.Second)

			tasks, err := ServiceClickup.ReturnTasks(lists.Lists[i].Id, 2)

			if err != nil {
				fmt.Println("Error return tasks: ", err.Error())
				return
			}

			for j := 0; j < len(tasks.Tasks); j++ {
				fmt.Println("Task ", j+1, "/", len(tasks.Tasks), " - ", tasks.Tasks[j].Name)

				auxEpicTaskId := ""

				if len(tasks.Tasks[j].LinkedTasks) == 0 {
					fmt.Println("Error 0 epics", " :: ", lists.Lists[i].Name, " :: ", tasks.Tasks[j].Name, " :: ", strings.ToLower(tasks.Tasks[j].Status.Status), " :: ", tasks.Tasks[j].Url)
					continue
				}

				if len(tasks.Tasks[j].LinkedTasks) > 1 {
					fmt.Println("Error 2 epics:", " :: ", lists.Lists[i].Name, " :: ", tasks.Tasks[j].Name, " :: ", strings.ToLower(tasks.Tasks[j].Status.Status), " :: ", tasks.Tasks[j].Url)
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

				task, err := ServiceClickup.ReturnTask(auxEpicTaskId)
				if err != nil {
					fmt.Println("Error return task: ", err.Error())
					return
				}

				if justVerify {
					time.Sleep(time.Second)
					ServiceClickup.VerifyTasks(task)
				} else {
					time.Sleep(time.Second)

					allSubTasksDone := true
					var timeSpent int64
					var requestTask TypeClickup.TaskRequest

					for k := 0; k < len(task.LinkedTasks); k++ {

						auxTaskId := ""

						if task.Id == task.LinkedTasks[k].LinkId {
							auxTaskId = task.LinkedTasks[k].TaskId
						} else {
							auxTaskId = task.LinkedTasks[k].LinkId
						}

						task, err := ServiceClickup.ReturnTask(auxTaskId)
						if err != nil {
							fmt.Println("Error task: " + err.Error())
							continue
						}
						auxDuoDate, _ := strconv.ParseInt(task.DueDate, 10, 64)
						auxStartDate, _ := strconv.ParseInt(task.StartDate, 10, 64)
						requestTask.TimeEstimate = requestTask.TimeEstimate + task.TimeEstimate
						timeSpent = timeSpent + task.TimeSpent
						if auxDuoDate > requestTask.DueDate {
							requestTask.DueDate = auxDuoDate
						}

						if auxStartDate != 0 && auxStartDate < requestTask.StartDate || requestTask.StartDate == 0 {
							requestTask.StartDate = auxStartDate
						}

						//caso tenha data no epic, não alterar qdo for 0
						if requestTask.StartDate == 0 && task.StartDate != "0" {
							requestTask.StartDate, _ = strconv.ParseInt(task.StartDate, 10, 64)
						}

						if task.Status.Status != "done" && task.Status.Status != "canceled" {
							allSubTasksDone = false
						}

						requestTask.Status, _ = ServiceClickup.RetNewStatus(task.Status.Status, task.Status.Status)
						task.Status.Status = requestTask.Status

						//precisa alterar o requirements
						//https://app.convisoappsec.com/scopes/553/projects/15296/project_requirements/232561

						//ServiceConvisoPlatform.ChangeActivitiesStatus(ServiceClickup.RetCustomFieldUrlConviso(task.CustomFields))

					}

					if allSubTasksDone {
						requestTask.Status = "done"
					}

					if (timeSpent - task.TimeSpent) > 0 {
						var taskTimeSpentRequest TypeClickup.TaskTimeSpentRequest
						taskTimeSpentRequest.Duration = timeSpent - task.TimeSpent
						taskTimeSpentRequest.Start = time.Now().UTC().UnixMilli()
						taskTimeSpentRequest.TaskId = task.Id
						ServiceClickup.RequestTaskTimeSpent(task.TeamId, taskTimeSpentRequest)
					}

					err := ServiceClickup.RequestPutTask(task.Id, requestTask)

					if err != nil {
						fmt.Println("Error task: " + err.Error())
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
		fmt.Println("3 - Count Deploys")
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
		case 3:
			ServiceConvisoPlatform.RetDeploys()
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
		VariablesConstant.CLICKUP_CUSTOMER_FIELD_ID,
		CustomFieldCustomerPosition,
	}

	customFieldUrlConvisoPlatform := TypeClickup.CustomFieldRequest{
		VariablesConstant.CLICKUP_URL_CONVISO_PLATFORM_FIELD_ID,
		"https://app.convisoappsec.com/scopes/" + strconv.Itoa(VariablesGlobal.Customer.PlatformID) + "/projects/" + project.Id,
	}

	customFieldTypeConsulting := TypeClickup.CustomFieldRequest{
		VariablesConstant.CLICKUP_TYPE_CONSULTING_FIELD_ID,
		strconv.Itoa(TypeEnumClickupTypeConsulting.STORE),
	}

	customFieldsMainTask := []TypeClickup.CustomFieldRequest{
		customFieldUrlConvisoPlatform,
		customFieldTypeConsulting,
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
				VariablesConstant.CLICKUP_URL_CONVISO_PLATFORM_FIELD_ID,
				convisoPlatformUrl.String(),
			}

			customFieldTypeConsultingSubTask := TypeClickup.CustomFieldRequest{
				VariablesConstant.CLICKUP_TYPE_CONSULTING_FIELD_ID,
				strconv.Itoa(TypeEnumClickupTypeConsulting.TASK),
			}

			customFieldsSubTask := []TypeClickup.CustomFieldRequest{
				customFieldUrlConvisoPlatformSubTask,
				customFieldTypeConsultingSubTask,
				customFieldCustomer}

			sanitizedHTMLTitle := ""
			sanitizedHTMLDescription := ""

			sanitizedHTMLTitle, err = html2text.FromString(project.Activities[i].Title)
			sanitizedHTMLDescription, err = html2text.FromString(project.Activities[i].Description)

			_, err := ServiceClickup.TaskCreateRequest(
				TypeClickup.TaskCreateRequest{
					sanitizedHTMLTitle,
					sanitizedHTMLDescription,
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
	// qdo atualizar uma história atualizar o projeto no conviso platform
	// qdo atualizar uma task atualizar também o requirements na platforma

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
