package main

import (
	"bytes"
	"flag"
	"fmt"
	"integration_platform_clickup_go/services/service_clickup"
	"integration_platform_clickup_go/services/service_conviso_platform"
	"integration_platform_clickup_go/types/type_clickup"
	"integration_platform_clickup_go/types/type_enum/enum_clickup_type_consulting"
	"integration_platform_clickup_go/types/type_platform"
	"integration_platform_clickup_go/utils/functions"
	"integration_platform_clickup_go/utils/variables_constant"
	"integration_platform_clickup_go/utils/variables_global"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jaytaylor/html2text"
	"golang.org/x/exp/slices"
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
	config := functions.LoadConfigsByYamlFile()

	fmt.Println("------Projets------")
	// Print the data
	for i := 0; i < len(config.Integrations); i++ {
		fmt.Println(i, " - ", config.Integrations[i].IntegrationName)
	}

	var input int
	fmt.Print("Enter the option: ")
	n, err := fmt.Scan(&input)
	if n < 1 || err != nil || input > len(config.Integrations)-1 {
		fmt.Println("Invalid Input")
		return
	}
	variables_global.Customer = config.Integrations[input]
}

func MenuSetupConfig() {
	var input int
	for ok := true; ok; ok = (input != 0) {
		fmt.Println("-----Menu Config-----")
		fmt.Println("Project Selected: ", variables_global.Customer.IntegrationName)
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

func VerifyErrorsProjectWithStore(list type_clickup.ListResponse) {
	VerifySubtask(list, int(enum_clickup_type_consulting.EPIC), int(enum_clickup_type_consulting.STORE))
	VerifySubtask(list, int(enum_clickup_type_consulting.STORE), int(enum_clickup_type_consulting.TASK))
	VerifyTasks(list)
}

func VerifyTasks(list type_clickup.ListResponse) {
	tasks, err := service_clickup.ReturnTasks(list.Id, int(enum_clickup_type_consulting.TASK))

	if err != nil {
		fmt.Println("Error VerifyTasks :: ", err.Error())
		return
	}

	for i := 0; i < len(tasks.Tasks); i++ {
		task, err := service_clickup.ReturnTask(tasks.Tasks[i].Id)

		if err != nil {
			fmt.Println("Error VerifyTasks GetTask :: ", err.Error())
			return
		}

		if task.Parent == "" {
			fmt.Println("TASK Without Store", " :: ", list.Name, " :: ", tasks.Tasks[i].Name, " :: ",
				strings.ToLower(tasks.Tasks[i].Status.Status), " :: ", tasks.Tasks[i].Url,
				" :: ", service_clickup.RetAssigness(tasks.Tasks[i].Assignees))
			continue
		}

		if strings.ToLower(task.Status.Status) != "backlog" && strings.ToLower(task.Status.Status) != "canceled" && strings.ToLower(task.Status.Status) != "blocked" {
			if task.DueDate == "" {
				fmt.Println("Task with errors: ", task.List.Name, " - ", task.Name, " - ", task.Name, " :: ", "DueDate empty", " :: ", task.Url,
					" :: ", service_clickup.RetAssigness(task.Assignees))
			}

			if task.StartDate == "" {
				fmt.Println("Task with errors: ", task.List.Name, " - ", task.Name, " - ", task.Name, " :: ", "StartDate empty", " :: ", task.Url,
					" :: ", service_clickup.RetAssigness(task.Assignees))
			}

			if task.TimeEstimate == 0 {
				fmt.Println("Task with errors: ", task.List.Name, " - ", task.Name, " - ", task.Name, " :: ", "TimeEstimate empty", " :: ", task.Url,
					" :: ", service_clickup.RetAssigness(task.Assignees))
			}

			if task.Status.Status == "done" && task.TimeSpent == 0 {
				fmt.Println("Task with errors: ", task.List.Name, " - ", task.Name, " - ", task.Name, " :: ", "TimeSpent empty", " :: ", task.Url,
					" :: ", service_clickup.RetAssigness(task.Assignees))
			}
		}
	}
}

func VerifySubtask(list type_clickup.ListResponse, customFieldTypeConsulting int, customFieldTypeConsultingSubTask int) {

	tasks, err := service_clickup.ReturnTasks(list.Id, customFieldTypeConsulting)

	if err != nil {
		fmt.Println("Error VerifySubtask :: ", err.Error())
		return
	}

	for i := 0; i < len(tasks.Tasks); i++ {
		task, err := service_clickup.ReturnTask(tasks.Tasks[i].Id)

		if err != nil {
			fmt.Println("Error VerifySubtask GetTask :: ", err.Error())
			return
		}

		if len(task.SubTasks) == 0 {
			fmt.Println(enum_clickup_type_consulting.ToString(customFieldTypeConsulting),
				" Without ",
				enum_clickup_type_consulting.ToString(customFieldTypeConsultingSubTask),
				" :: ", list.Name, " :: ", tasks.Tasks[i].Name, " :: ", strings.ToLower(tasks.Tasks[i].Status.Status), " :: ", tasks.Tasks[i].Url,
				" :: ", service_clickup.RetAssigness(tasks.Tasks[i].Assignees))
			continue
		}

		for j := 0; j < len(task.SubTasks); j++ {
			subTask, err := service_clickup.ReturnTask(task.SubTasks[j].Id)
			if err != nil {
				fmt.Println("Error VerifySubtask GetTask GetSubTask :: ", err.Error())
				return
			}

			customFieldsSubTask := service_clickup.RetCustomFieldTypeConsulting(subTask.CustomFields)

			if customFieldsSubTask != customFieldTypeConsultingSubTask {
				fmt.Println(
					subTask.Name,
					" should be ",
					enum_clickup_type_consulting.ToString(customFieldTypeConsultingSubTask),
					" but is ",
					enum_clickup_type_consulting.ToString(customFieldsSubTask),
					" :: ", list.Name, " :: ", strings.ToLower(subTask.Status.Status), " :: ", subTask.Url,
					" :: ", service_clickup.RetAssigness(subTask.Assignees))
				continue
			}
		}
	}
}

func UpdateProjectWithStore(list type_clickup.ListResponse) {
	UpdateSubtask(list, enum_clickup_type_consulting.TASK, enum_clickup_type_consulting.STORE)
	UpdateSubtask(list, enum_clickup_type_consulting.STORE, enum_clickup_type_consulting.EPIC)
}

func UpdateSubtask(list type_clickup.ListResponse, typeConsultingTask int, typeConsultingParent int) {

	tasks, err := service_clickup.ReturnTasks(list.Id, typeConsultingTask)

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

		taskParent, err := service_clickup.ReturnTask((tasks.Tasks[i].Parent))

		var convisoPlatformProject type_platform.Project

		if err != nil {
			fmt.Println("Error UpdateSubtask GetTask Parent :: ", err.Error())
			continue
		}

		var requestTask type_clickup.TaskRequestStore
		requestTask.Status = taskParent.Status.Status
		requestTask.DueDate, _ = strconv.ParseInt(taskParent.DueDate, 10, 64)
		requestTask.StartDate, _ = strconv.ParseInt(taskParent.StartDate, 10, 64)
		allTaskDone := true
		hasUpdate := false
		for j := 0; j < len(taskParent.SubTasks); j++ {
			subTask, err := service_clickup.ReturnTask(taskParent.SubTasks[j].Id)
			if err != nil {
				fmt.Println("Error UpdateSubtask GetTask SubTask :: ", err.Error())
				return
			}
			var auxStartDate int64
			var auxDueDate int64

			auxStartDate, _ = strconv.ParseInt(subTask.StartDate, 10, 64)
			auxDueDate, _ = strconv.ParseInt(subTask.DueDate, 10, 64)
			if (auxStartDate < requestTask.StartDate && auxStartDate != 0) || (auxStartDate != 0 && requestTask.StartDate == 0) {
				requestTask.StartDate = auxStartDate
				hasUpdate = true
			}

			if (auxDueDate > requestTask.DueDate && auxDueDate != 0) || (auxDueDate != 0 && requestTask.DueDate == 0) {
				requestTask.DueDate = auxDueDate
				hasUpdate = true
			}

			hasUpdateStatus := false
			requestTask.Status, hasUpdateStatus = service_clickup.RetNewStatus(requestTask.Status, subTask.Status.Status)

			if hasUpdateStatus {
				hasUpdate = true
			}

			if subTask.Status.Status != "done" && subTask.Status.Status != "canceled" {
				allTaskDone = false
			}

			if taskParent.CustomField.TypeConsulting == enum_clickup_type_consulting.STORE &&
				taskParent.CustomField.LinkConvisoPlatform != "" && convisoPlatformProject.Id == "" {

				projectId, err := service_conviso_platform.RetProjectIdCustomField(taskParent.CustomField.LinkConvisoPlatform)

				if err == nil {
					convisoPlatformProject, err = service_conviso_platform.GetProject(projectId)
					if err != nil {
						fmt.Println("Error GetProject Conviso Platform :: ", err.Error())
					}
				} else {
					fmt.Println("Error RetProjectIdCustomField Conviso Platform :: ", err.Error())
				}
			}

			//update the activity in conviso platform project
			err = service_conviso_platform.UpdateActivityRequirement(subTask, convisoPlatformProject)

			if err != nil {
				fmt.Println("Task ", subTask.Name, " not possible update requirement activity in Conviso Platform")
			}

		}

		if allTaskDone {
			requestTask.Status = "done"
			hasUpdate = true
		}

		if hasUpdate {
			err = service_clickup.RequestPutTaskStore(taskParent.Id, requestTask)
			if err != nil {
				fmt.Println("Store not possible update in clickup")
			} else {
				err = service_conviso_platform.UpdateProjectRest(requestTask, convisoPlatformProject.Id, taskParent.TimeEstimate)
				if err != nil {
					fmt.Println("Store not possible update in conviso platform")
				}
			}
		}
	}
}

func UpdateClickUpConvisoPlatform(justVerify bool) {
	fmt.Println("...Starting ClickUp Automation...")

	for i := 0; i < len(variables_global.Config.Integrations); i++ {

		fmt.Println("Found List ", variables_global.Config.Integrations[i].IntegrationName)

		list, error := service_clickup.ReturnList(variables_global.Config.Integrations[i].ClickUpListId)

		if error != nil {
			fmt.Println("Error loading list ", variables_global.Config.Integrations[i].IntegrationName)
			continue
		}

		variables_global.Customer = variables_global.Config.Integrations[i]

		if justVerify {
			VerifyErrorsProjectWithStore(list)
			//return
		}

		if !justVerify {
			UpdateProjectWithStore(list)
			//return
		}
	}

	fmt.Println("...Finishing ClickUp Automation...")
}

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
			service_conviso_platform.InputSearchRequimentsPlatform()
		case 2:
			service_conviso_platform.InputSearchProjectTypesPlatform()
		case 3:
			service_conviso_platform.RetDeploys()
		default:
			fmt.Println("Invalid Input")
		}
	}
}

func MainMenu() {
	var input int

	for ok := true; ok; ok = (input != 0) {
		fmt.Println("-----Main Menu-----")
		fmt.Println("App Version: ", variables_constant.VERSION)
		fmt.Println("SO: ", runtime.GOOS)
		fmt.Println("Arch: ", runtime.GOARCH)
		fmt.Println("Project Selected: ", variables_global.Customer.IntegrationName)
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
		case 1:
			MenuSetupConfig()
		case 2:
			if variables_global.Customer.PlatformID == 0 {
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
	requirementIds := ""
	typeId := 10
	SubTaskReqActivies := "n"

	//fmt.Print("Label: ")
	label := functions.GetTextWithSpace("Label: ")

	//fmt.Print("Goal: ")
	goal := functions.GetTextWithSpace("Goal: ")

	//fmt.Print("Scope: ")
	scope := functions.GetTextWithSpace("Scope: ")

	fmt.Print("TypeId (Consulting = 10): ")
	n, err := fmt.Scan(&typeId)
	if n < 1 || err != nil {
		fmt.Println("Invalid Input")
		return
	}

	fmt.Print("Requirement ID (1;2;3): ")
	n, err = fmt.Scan(&requirementIds)
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

	createConvisoPlatform := type_platform.ProjectCreateInputRequest{variables_global.Customer.PlatformID,
		label, goal, scope, typeId,
		functions.ConvertStringToArrayInt(requirementIds),
		time.Now().Add(time.Hour * 24).Format("2006-01-02"), "1"}

	err = service_conviso_platform.AddPlatformProject(createConvisoPlatform)

	if err != nil {
		fmt.Println("Error CreateProject: ", err.Error())
	}

	project, err := service_conviso_platform.ConfirmProjectCreate(variables_global.Customer.PlatformID, label)

	if err != nil {
		fmt.Println("Erro CreateProject: Contact the system administrator")
		return
	}

	// CustomFieldCustomerPosition, err := service_clickup.RetCustomerPosition()
	// if err != nil {
	// 	fmt.Println("Error CreateProject: CustomerCustomField error...")
	// 	CustomFieldCustomerPosition = ""
	// }

	// customFieldCustomer := type_clickup.CustomFieldRequest{
	// 	variables_constant.CLICKUP_CUSTOMER_FIELD_ID,
	// 	CustomFieldCustomerPosition,
	// }

	customFieldUrlConvisoPlatform := type_clickup.CustomFieldRequest{
		variables_constant.CLICKUP_URL_CONVISO_PLATFORM_FIELD_ID,
		"https://app.convisoappsec.com/scopes/" + strconv.Itoa(variables_global.Customer.PlatformID) + "/projects/" + project.Id,
	}

	customFieldTypeConsulting := type_clickup.CustomFieldRequest{
		variables_constant.CLICKUP_TYPE_CONSULTING_FIELD_ID,
		strconv.Itoa(enum_clickup_type_consulting.STORE),
	}

	customFieldsMainTask := []type_clickup.CustomFieldRequest{
		customFieldUrlConvisoPlatform,
		customFieldTypeConsulting,
		//customFieldCustomer,
	}

	assignessTask := []int64{variables_global.Config.ConfclickUp.User}

	//create main
	taskMainClickup, err := service_clickup.TaskCreateRequest(
		type_clickup.TaskCreateRequest{
			project.Label,
			project.Scope,
			"backlog",
			true,
			"",
			"",
			customFieldsMainTask,
			assignessTask,
		})

	if err != nil {
		fmt.Println("Error CreateProject: Problem create ClickUpTask :: ", err.Error())
		return
	}

	if strings.ToLower(SubTaskReqActivies) == "y" {

		for i := 0; i < len(project.Activities); i++ {
			var convisoPlatformUrl bytes.Buffer
			convisoPlatformUrl.WriteString(variables_constant.CONVISO_PLATFORM_URL_BASE)
			convisoPlatformUrl.WriteString("scopes/")
			convisoPlatformUrl.WriteString(strconv.Itoa(variables_global.Customer.PlatformID))
			convisoPlatformUrl.WriteString("/projects/")
			convisoPlatformUrl.WriteString(project.Id)
			convisoPlatformUrl.WriteString("/project_requirements/")
			convisoPlatformUrl.WriteString(project.Activities[i].Id)

			customFieldUrlConvisoPlatformSubTask := type_clickup.CustomFieldRequest{
				variables_constant.CLICKUP_URL_CONVISO_PLATFORM_FIELD_ID,
				convisoPlatformUrl.String(),
			}

			customFieldTypeConsultingSubTask := type_clickup.CustomFieldRequest{
				variables_constant.CLICKUP_TYPE_CONSULTING_FIELD_ID,
				strconv.Itoa(enum_clickup_type_consulting.TASK),
			}

			customFieldsSubTask := []type_clickup.CustomFieldRequest{
				customFieldUrlConvisoPlatformSubTask,
				customFieldTypeConsultingSubTask,
				// customFieldCustomer,
			}

			sanitizedHTMLTitle := ""
			sanitizedHTMLDescription := ""

			sanitizedHTMLTitle, err = html2text.FromString(project.Activities[i].Title)
			sanitizedHTMLDescription, err = html2text.FromString(project.Activities[i].Description)

			_, err := service_clickup.TaskCreateRequest(
				type_clickup.TaskCreateRequest{
					sanitizedHTMLTitle,
					sanitizedHTMLDescription,
					"backlog",
					true,
					taskMainClickup.Id,
					taskMainClickup.Id,
					customFieldsSubTask,
					assignessTask,
				})
			if err != nil {
				fmt.Println("Error CreateProject: Problem create ClickUp SubTask ", project.Activities[i].Title)
			}
		}
	}

	fmt.Println("Create Task Success!")
}

func main() {

	integrationJustVerify := flag.Bool("iv", false, "Verify if clickup tasks is ok")
	integrationUpdateTasks := flag.Bool("iu", false, "Update Conviso Platform and ClickUp Tasks")
	deploy := flag.Bool("d", false, "See info about deploys")
	version := flag.Bool("v", false, "Script Version")

	variables_global.Config = functions.LoadConfigsByYamlFile()

	if variables_global.Config.ConfigGeneral.IntegrationDefault != -1 {
		variables_global.Customer = variables_global.Config.Integrations[variables_global.Config.ConfigGeneral.IntegrationDefault]
	}

	flag.Parse()

	if *integrationJustVerify {
		UpdateClickUpConvisoPlatform(true)
		os.Exit(0)
	}

	if *integrationUpdateTasks {
		UpdateClickUpConvisoPlatform(false)
		os.Exit(0)
	}

	if *deploy {
		service_conviso_platform.RetDeploys()
		os.Exit(0)
	}

	if *version {
		fmt.Println("Script Version: ", variables_constant.VERSION)
		os.Exit(0)
	}

	MainMenu()

}
