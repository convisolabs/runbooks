package main

import (
	"bytes"
	"flag"
	"fmt"
	crawler_service "integration_platform_clickup_go/services/crawler"
	"integration_platform_clickup_go/services/service_clickup"
	"integration_platform_clickup_go/services/service_conviso_platform"
	slack_service "integration_platform_clickup_go/services/slack"
	"integration_platform_clickup_go/types/type_clickup"
	"integration_platform_clickup_go/types/type_config"
	"integration_platform_clickup_go/types/type_enum/enum_clickup_ps_team"
	"integration_platform_clickup_go/types/type_enum/enum_clickup_statuses"
	"integration_platform_clickup_go/types/type_enum/enum_clickup_type_ps_hierarchy"
	"integration_platform_clickup_go/types/type_enum/enum_main_action"
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
	"gopkg.in/yaml.v2"
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
	config := variables_global.Config

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
		case 1:
			LoadProjects()
		default:
			fmt.Println("Invalid Input")
		}
	}
}

func VerifyErrorsProjectWithStore(list type_clickup.ListResponse) {
	VerifySubtask(list, int(enum_clickup_type_ps_hierarchy.EPIC), int(enum_clickup_type_ps_hierarchy.STORE))
	VerifySubtask(list, int(enum_clickup_type_ps_hierarchy.STORE), int(enum_clickup_type_ps_hierarchy.TASK))
	VerifyTasks(list)
}

func VerifyTasks(list type_clickup.ListResponse) {

	page := 0

	for {

		tasks, err := service_clickup.ReturnTasks(list.Id,
			type_clickup.SearchTask{
				TaskType:      int(enum_clickup_type_ps_hierarchy.TASK),
				Page:          page,
				DateUpdatedGt: time.Now().Add(-time.Hour * 240).UTC().UnixMilli(),
				IncludeClosed: false,
				SubTasks:      true,
				TaskStatuses:  "",
			},
		)

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

			if strings.ToLower(task.Status.Status) != "backlog" &&
				strings.ToLower(task.Status.Status) != "closed" &&
				!service_clickup.CheckSpecificTag(task.Tags, "não executada") {

				if task.Parent == "" {
					fmt.Println("TASK Without Store", " :: ", variables_global.Customer.IntegrationName, " :: ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", task.Url,
						" :: ", service_clickup.RetAssigness(task.Assignees))
					continue
				}

				if task.DueDate == "" {
					fmt.Println("Task with errors: ", variables_global.Customer.IntegrationName, " - ", task.Name, " - ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", "DueDate empty", " :: ", task.Url,
						" :: ", service_clickup.RetAssigness(task.Assignees))
				}

				if task.StartDate == "" {
					fmt.Println("Task with errors: ", variables_global.Customer.IntegrationName, " - ", task.Name, " - ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", "StartDate empty", " :: ", task.Url,
						" :: ", service_clickup.RetAssigness(task.Assignees))
				}

				if task.TimeEstimate == 0 {
					fmt.Println("Task with errors: ", variables_global.Customer.IntegrationName, " - ", task.Name, " - ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", "TimeEstimate empty", " :: ", task.Url,
						" :: ", service_clickup.RetAssigness(task.Assignees))
				}

				if task.Status.Status == "done" && task.TimeSpent == 0 {
					fmt.Println("Task with errors: ", variables_global.Customer.IntegrationName, " - ", task.Name, " - ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", "TimeSpent empty", " :: ", task.Url,
						" :: ", service_clickup.RetAssigness(task.Assignees))
				}

				if len(task.CustomField.PSTeam) == 0 {
					fmt.Println("Task with errors: ", variables_global.Customer.IntegrationName, " - ", task.Name, " - ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", "PS-Team empty", " :: ", task.Url,
						" :: ", service_clickup.RetAssigness(task.Assignees))
				}

				if variables_global.Customer.ValidatePSCustomer && len(task.CustomField.PSCustomer) == 0 {
					fmt.Println("Task with errors: ", variables_global.Customer.IntegrationName, " - ", task.Name, " - ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", "PS-Customer empty", " :: ", task.Url,
						" :: ", service_clickup.RetAssigness(task.Assignees))
				}
			}
		}

		if tasks.LastPage {
			break
		}

		page++
	}
}

func VerifySubtask(list type_clickup.ListResponse, customFieldTypeConsulting int, customFieldTypeConsultingSubTask int) {

	page := 0

	for {

		tasks, err := service_clickup.ReturnTasks(list.Id,
			type_clickup.SearchTask{
				TaskType:      customFieldTypeConsulting,
				Page:          page,
				DateUpdatedGt: time.Now().Add(-time.Hour * 240).UTC().UnixMilli(),
				IncludeClosed: false,
				SubTasks:      true,
				TaskStatuses:  "",
			},
		)

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

			if strings.EqualFold(task.Parent, "") && customFieldTypeConsulting != int(enum_clickup_type_ps_hierarchy.EPIC) {
				fmt.Println("Store  Without EPIC",
					" :: ", variables_global.Customer.IntegrationName, " :: ", task.Name,
					" :: ", strings.ToLower(task.Status.Status), " :: ", task.Url,
					" :: ", service_clickup.RetAssigness(task.Assignees))
				continue
			}

			if len(task.SubTasks) == 0 {
				fmt.Println(enum_clickup_type_ps_hierarchy.ToString(customFieldTypeConsulting),
					" Without ",
					enum_clickup_type_ps_hierarchy.ToString(customFieldTypeConsultingSubTask),
					" :: ", variables_global.Customer.IntegrationName, " :: ", task.Name,
					" :: ", strings.ToLower(task.Status.Status), " :: ", task.Url,
					" :: ", service_clickup.RetAssigness(task.Assignees))
				continue
			}

			if len(task.CustomField.PSTeam) == 0 {
				fmt.Println("EPIC or Story without PS-TEAM: ", variables_global.Customer.IntegrationName, " :: ", task.Name, " :: ",
					strings.ToLower(task.Status.Status), " :: ", task.Url,
					" :: ", service_clickup.RetAssigness(task.Assignees))
			}

			if variables_global.Customer.ValidatePSCustomer && len(task.CustomField.PSCustomer) == 0 {
				fmt.Println("EPIC or Story without PS-Customer: ", variables_global.Customer.IntegrationName, " :: ", task.Name, " :: ",
					strings.ToLower(task.Status.Status), " :: ", task.Url,
					" :: ", service_clickup.RetAssigness(task.Assignees))
			}

			if customFieldTypeConsulting == int(enum_clickup_type_ps_hierarchy.STORE) && variables_global.Customer.ValidateTag {
				if !service_clickup.CheckTags(task.Tags) {
					fmt.Println("Story without TAGS", " :: ", variables_global.Customer.IntegrationName, " :: ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", task.Url,
						" :: ", service_clickup.RetAssigness(task.Assignees))

				}

				if variables_global.Customer.ValidatePSConvisoPlatformLink && (task.CustomField.PSConvisoPlatformLink == "" || !strings.Contains(task.CustomField.PSConvisoPlatformLink, "/projects/")) {
					fmt.Println("Story without Conviso Platform URL: ", " :: ", variables_global.Customer.IntegrationName,
						" :: ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", task.Url,
						" :: ", service_clickup.RetAssigness(task.Assignees))
				}
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
						enum_clickup_type_ps_hierarchy.ToString(customFieldTypeConsultingSubTask),
						" but is ",
						enum_clickup_type_ps_hierarchy.ToString(customFieldsSubTask),
						" :: ", variables_global.Customer.IntegrationName, " :: ",
						subTask.Name, " :: ",
						strings.ToLower(subTask.Status.Status),
						" :: ", subTask.Url, " :: ", service_clickup.RetAssigness(subTask.Assignees))
				}
			}
		}

		if tasks.LastPage {
			break
		}

		page++
	}

}

func ListStoryInProgress(list type_clickup.ListResponse) {
	page := 0

	for {
		tasks, err := service_clickup.ReturnTasks(list.Id,
			type_clickup.SearchTask{
				TaskType:      enum_clickup_type_ps_hierarchy.STORE,
				Page:          page,
				DateUpdatedGt: 0,
				IncludeClosed: false,
				SubTasks:      true,
				TaskStatuses:  enum_clickup_statuses.IN_PROGRESS,
			},
		)

		if err != nil {
			fmt.Println("Error ListStoryInProgress :: ", err.Error())
			return
		}

		for i := 0; i < len(tasks.Tasks); i++ {

			var dtStart time.Time
			var dtDuoDate time.Time

			dtIntAux, err := strconv.ParseInt(tasks.Tasks[i].StartDate, 10, 64)
			if err == nil {
				dtStart = time.UnixMilli(dtIntAux)
			}

			dtIntAux, err = strconv.ParseInt(tasks.Tasks[i].DueDate, 10, 64)
			if err == nil {
				dtDuoDate = time.UnixMilli(dtIntAux)
			}

			fmt.Println("Story In Progress",
				";", variables_global.Customer.IntegrationName,
				";", tasks.Tasks[i].Name,
				";", tasks.Tasks[i].Url,
				";", dtStart.Format("02/01/2006"),
				";", dtDuoDate.Format("02/01/2006"),
				";", service_clickup.RetAssigness(tasks.Tasks[i].Assignees))
		}

		if tasks.LastPage {
			break
		}

		page++
	}
}

func ListTasksInClosed(list type_clickup.ListResponse) {
	ListTasksInClosedByPSHierarchy(list, enum_clickup_type_ps_hierarchy.EPIC)
	ListTasksInClosedByPSHierarchy(list, enum_clickup_type_ps_hierarchy.STORE)
	ListTasksInClosedByPSHierarchy(list, enum_clickup_type_ps_hierarchy.TASK)
}

func UpdateTasksInDoneToClosed(list type_clickup.ListResponse) {
	UpdateTasksInDoneToClosedPSHierarchy(list, enum_clickup_type_ps_hierarchy.TASK)
	UpdateTasksInDoneToClosedPSHierarchy(list, enum_clickup_type_ps_hierarchy.STORE)
	UpdateTasksInDoneToClosedPSHierarchy(list, enum_clickup_type_ps_hierarchy.EPIC)
}

func UpdateTasksInDoneToClosedPSHierarchy(list type_clickup.ListResponse, psHierarchy int) {
	page := 0

	for {
		tasks, err := service_clickup.ReturnTasks(list.Id,
			type_clickup.SearchTask{
				TaskType:      psHierarchy,
				Page:          page,
				DateUpdatedGt: 0,
				IncludeClosed: false,
				SubTasks:      true,
				TaskStatuses:  "done",
			},
		)

		if err != nil {
			fmt.Println("Error UpdateTasksInDoneToClosedPSHierarchy :: ", err.Error())
			return
		}

		for i := 0; i < len(tasks.Tasks); i++ {
			err = service_clickup.RequestPutTaskStatus(tasks.Tasks[i].Id, type_clickup.TaskRequestStatus{
				Status: "closed",
			})

			if err != nil {
				fmt.Println("Error UpdateTasksInDoneToClosedPSHierarchy :: ", tasks.Tasks[i].Url, " :: ", err.Error())
				return
			}
		}

		if tasks.LastPage {
			break
		}

		page++
	}
}

func ListTasksInClosedByPSHierarchy(list type_clickup.ListResponse, psHierarchy int) {
	page := 0

	for {
		tasks, err := service_clickup.ReturnTasks(list.Id,
			type_clickup.SearchTask{
				TaskType:      psHierarchy,
				Page:          page,
				DateUpdatedGt: time.Now().Add(-time.Hour * 60).UTC().UnixMilli(),
				IncludeClosed: true,
				SubTasks:      true,
				TaskStatuses:  "closed",
			},
		)

		if err != nil {
			fmt.Println("Error ListTasksInClosedByPSHierarchy :: ", err.Error())
			return
		}

		for i := 0; i < len(tasks.Tasks); i++ {
			fmt.Println(enum_clickup_type_ps_hierarchy.ToString(psHierarchy), " Closed ",
				" :: ", variables_global.Customer.IntegrationName,
				" :: ", tasks.Tasks[i].Name,
				" :: ", tasks.Tasks[i].Url,
				" :: ", service_clickup.RetAssigness(tasks.Tasks[i].Assignees))
		}

		if tasks.LastPage {
			break
		}

		page++
	}
}

func UpdateProjectWithStore(list type_clickup.ListResponse) {
	UpdateTask(list, enum_clickup_type_ps_hierarchy.TASK, enum_clickup_type_ps_hierarchy.STORE)
	UpdateTask(list, enum_clickup_type_ps_hierarchy.STORE, enum_clickup_type_ps_hierarchy.EPIC)
}

func UpdateTask(list type_clickup.ListResponse, typeConsultingTask int, typeConsultingParent int) {
	page := 0

	for {

		tasks, err := service_clickup.ReturnTasks(list.Id,
			type_clickup.SearchTask{
				TaskType:      typeConsultingTask,
				Page:          page,
				DateUpdatedGt: time.Now().Add(-time.Hour * 240).UTC().UnixMilli(),
				IncludeClosed: false,
				SubTasks:      true,
				TaskStatuses:  "",
			},
		)

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

			if taskParent.CustomField.PSProjectHierarchy != typeConsultingParent {
				fmt.Println(
					taskParent.Id, " :: ", taskParent.Name, " :: ", taskParent.Url,
					" :: ", " isn't type ", enum_clickup_type_ps_hierarchy.ToString(typeConsultingParent),
				)
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

				if !strings.EqualFold(subTask.Status.Status, "done") &&
					!strings.EqualFold(subTask.Status.Status, "canceled") &&
					!strings.EqualFold(subTask.Status.Status, "closed") {
					allTaskDone = false
				}

				if taskParent.CustomField.PSProjectHierarchy == enum_clickup_type_ps_hierarchy.STORE &&
					taskParent.CustomField.PSConvisoPlatformLink != "" && convisoPlatformProject.Id == "" {

					projectId, err := service_conviso_platform.RetProjectIdCustomField(taskParent.CustomField.PSConvisoPlatformLink)

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

			if taskParent.CustomField.PSProjectHierarchy == enum_clickup_type_ps_hierarchy.STORE && variables_global.Customer.ValidateTag {
				deliveryPoint := service_clickup.RetDeliveryPointTag(taskParent.Tags)
				deliveruPointString := strconv.Itoa(deliveryPoint)
				if deliveryPoint != 0 && !strings.EqualFold(deliveruPointString, taskParent.CustomField.PSDeliveryPoints) {

					err = service_clickup.RequestSetValueCustomField(taskParent.Id,
						variables_constant.CLICKUP_CUSTOM_FIELD_PS_DELIVERY_POINTS,
						type_clickup.CustomFieldValueRequest{
							deliveruPointString,
						},
					)

					if err != nil {
						fmt.Println("Store not possible update delivery points")
					}

				}
			}
		}

		if tasks.LastPage {
			break
		}

		page++
	}
}

func MainAction(mainAction int) {
	fmt.Println("...Starting ClickUp Automation...")

	for i := 0; i < len(variables_global.Config.Integrations); i++ {

		fmt.Println("Found List ", variables_global.Config.Integrations[i].IntegrationName)
		fmt.Println("Begin: ", time.Now().Format("2006-01-02 15:04:05"))

		list, error := service_clickup.ReturnList(variables_global.Config.Integrations[i].ClickUpListId)

		if error != nil {
			fmt.Println("Error loading list ", variables_global.Config.Integrations[i].IntegrationName)
			continue
		}

		variables_global.Customer = variables_global.Config.Integrations[i]

		switch mainAction {
		case enum_main_action.TASKS_VERIFY:
			VerifyErrorsProjectWithStore(list)

		case enum_main_action.TASKS_UPDATE:
			UpdateProjectWithStore(list)

		case enum_main_action.TASKS_INPROGRESS:
			ListStoryInProgress(list)

		case enum_main_action.TASKS_INCLOSED:
			ListTasksInClosed(list)

		case enum_main_action.ASSETS_NEW_CP_FORTIFY:
			if variables_global.Customer.AssetNewFortify {
				AssetsNew(variables_global.Customer)
			}
		case enum_main_action.TASKS_UPDATE_DONE_CLOSED:
			UpdateTasksInDoneToClosed(list)

		}

		fmt.Println("Finish: ", time.Now().Format("2006-01-02 15:04:05"))
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
		fmt.Println("4 - Save Yaml File Test")

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
		case 4:
			//SaveYamlFileTest()
			AssetsNew(variables_global.Customer)
		default:
			fmt.Println("Invalid Input")
		}
	}
}

func SaveYamlFileTest() {
	fileName := "test.yaml"
	yamlData, err := yaml.Marshal(&variables_global.Config)
	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	err = os.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		panic("Unable to write data into the file")
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

	customFieldUrlConvisoPlatform := type_clickup.CustomFieldRequest{
		variables_constant.CLICKUP_CUSTOM_FIELD_PS_CP_LINK,
		"https://app.convisoappsec.com/scopes/" + strconv.Itoa(variables_global.Customer.PlatformID) + "/projects/" + project.Id,
	}

	customFieldPSHierarchy := type_clickup.CustomFieldRequest{
		variables_constant.CLICKUP_CUSTOM_FIELD_PS_HIERARCHY,
		strconv.Itoa(enum_clickup_type_ps_hierarchy.STORE),
	}

	customFieldPSTeam := type_clickup.CustomFieldRequest{
		variables_constant.CLICKUP_CUSTOM_FIELD_PS_TEAM_ID,
		strconv.Itoa(enum_clickup_ps_team.CONSULTING),
	}

	customerOrder, err := service_clickup.RetClickUpDropDownPosition(variables_global.Customer.ClickUpListId, variables_constant.CLICKUP_CUSTOM_FIELD_PS_CUSTOMER_ID,
		variables_global.Customer.ClickUpCustomerList)

	if err != nil {
		fmt.Println("Error customerOrder: Contact the system administrator")
		return
	}

	customFieldPSCustomer := type_clickup.CustomFieldRequest{
		variables_constant.CLICKUP_CUSTOM_FIELD_PS_CUSTOMER_ID,
		strconv.Itoa(customerOrder),
	}

	customFieldsMainTask := []type_clickup.CustomFieldRequest{
		customFieldUrlConvisoPlatform,
		customFieldPSHierarchy,
		customFieldPSTeam,
		customFieldPSCustomer,
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
				variables_constant.CLICKUP_CUSTOM_FIELD_PS_CP_LINK,
				convisoPlatformUrl.String(),
			}

			customFieldTypeConsultingSubTask := type_clickup.CustomFieldRequest{
				variables_constant.CLICKUP_CUSTOM_FIELD_PS_HIERARCHY,
				strconv.Itoa(enum_clickup_type_ps_hierarchy.TASK),
			}

			customFieldsSubTask := []type_clickup.CustomFieldRequest{
				customFieldUrlConvisoPlatformSubTask,
				customFieldTypeConsultingSubTask,
				customFieldPSTeam,
				customFieldPSCustomer,
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

func AssetsNew(integration type_config.ConfigTypeIntegration) {

	var urlBase bytes.Buffer

	urlBase.WriteString(variables_constant.CONVISO_PLATFORM_URL_BASE)
	urlBase.WriteString("scopes/")
	urlBase.WriteString(strconv.Itoa(integration.PlatformID))
	urlBase.WriteString("/integrations/fortify/select_projects?page={1}")

	page := 1

	slackService := slack_service.SlackServiceNew()
	crawlerService := crawler_service.CrawlerServiceNew(slackService)

	for {
		urlPage := strings.Replace(urlBase.String(), "{1}", strconv.Itoa(page), -1)
		cont := crawlerService.Exec(integration.PlatformID, urlPage)
		if !cont {
			break
		}
		page++
	}
}

func InitialCheck() bool {
	ret := true

	err := error(nil)

	variables_global.Config, err = functions.LoadConfigsByYamlFile()

	if err != nil {
		fmt.Println("YAML File Problem", variables_constant.CLICKUP_TOKEN_NAME, " is empty!")
		ret = false
	}

	if os.Getenv(variables_constant.CLICKUP_TOKEN_NAME) == "" {
		fmt.Println("Variable ", variables_constant.CLICKUP_TOKEN_NAME, " is empty!")
		ret = false
	}

	if os.Getenv(variables_constant.CONVISO_PLATFORM_TOKEN_NAME) == "" {
		fmt.Println("Variable ", variables_constant.CONVISO_PLATFORM_TOKEN_NAME, " is empty!")
		ret = false
	}

	return ret
}

func SetDefaultValue() {
	if variables_global.Config.ConfclickUp.HttpAttempt == nil {
		httpAttempt := 3
		variables_global.Config.ConfclickUp.HttpAttempt = &httpAttempt
	}
}

func main() {
	/*
		TODO LIST
			remover gambiarra verificar CP criou projeto
			separar as atualizações do cp e clickup, hoje tem uma variável, has update, mas deveria ter algo do tipo hasupdate cp e hasupcate clickup
			Verificar o status do conviso platform x clickup e atualizar o CP exemplo está in progress no Clickup e no CP planned https://app.convisoappsec.com/spa/scopes/413/projects/19138?locale=en
			qdo não encontrar um cliente no campo PS Customer, não quebrar a aplicação, selecionar o primeiro da lista ou algo assim
			verificar possibilidade de melhorar a função de recuperar o customfield do clickup na função returntask
	*/

	AssetsNew(variables_global.Customer)

	if !InitialCheck() {
		fmt.Println("You need to correct the above information before rerunning the application")
		fmt.Println("Press the Enter Key to finish!")
		fmt.Scanln()
		os.Exit(0)
	}

	SetDefaultValue()

	integrationJustVerify := flag.Bool("tv", false, "Verify if clickup tasks is ok")
	integrationUpdateTasks := flag.Bool("tu", false, "Update Conviso Platform and ClickUp Tasks")
	integrationListTasksInProgress := flag.Bool("tsip", false, "List Clickup Stories In Progress")
	integrationListTasksClosed := flag.Bool("tsd", false, "List Clickup Epics, Stories and Tasks in Closed")
	crawlerAssetNewCP := flag.Bool("can", false, "Search New Assets Fortify Integration Conviso Platform")
	integrationUpdateTasksDone := flag.Bool("tud", false, "Change tasks done to closed")
	deploy := flag.Bool("d", false, "See info about deploys")
	version := flag.Bool("v", false, "Script Version")

	if variables_global.Config.ConfigGeneral.IntegrationDefault != -1 {
		if len(variables_global.Config.Integrations) <= variables_global.Config.ConfigGeneral.IntegrationDefault {
			variables_global.Config.ConfigGeneral.IntegrationDefault = 0
		}
		variables_global.Customer = variables_global.Config.Integrations[variables_global.Config.ConfigGeneral.IntegrationDefault]
	}

	flag.Parse()

	if *integrationJustVerify {
		MainAction(enum_main_action.TASKS_VERIFY)
		os.Exit(0)
	}

	if *integrationUpdateTasks {
		MainAction(enum_main_action.TASKS_UPDATE)
		os.Exit(0)
	}

	if *integrationListTasksInProgress {
		MainAction(enum_main_action.TASKS_INPROGRESS)
		os.Exit(0)
	}

	if *integrationListTasksClosed {
		MainAction(enum_main_action.TASKS_INCLOSED)
		os.Exit(0)
	}

	if *crawlerAssetNewCP {
		MainAction(enum_main_action.ASSETS_NEW_CP_FORTIFY)
		os.Exit(0)
	}

	if *integrationUpdateTasksDone {
		MainAction(enum_main_action.TASKS_UPDATE_DONE_CLOSED)
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
