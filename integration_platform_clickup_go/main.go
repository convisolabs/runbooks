package main

import (
	"bytes"
	"flag"
	"fmt"
	"integration_platform_clickup_go/services/clickup_service"
	cp_service "integration_platform_clickup_go/services/cp"
	"integration_platform_clickup_go/types/type_clickup"
	"integration_platform_clickup_go/types/type_enum/enum_clickup_ps_team"
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
)

const BANNER = `
____  _       _    __                       ____ _ _      _    _   _       
|  _ \| | __ _| |_ / _| ___  _ __ _ __ ___  / ___| (_) ___| | _| | | |_ __  
| |_) | |/ _∎ | __| |_ / _ \| '__| '_ ∎ _ \| |   | | |/ __| |/ / | | | '_ \ 
|  __/| | (_| | |_|  _| (_) | |  | | | | | | |___| | | (__|   <| |_| | |_) |
|_|   |_|\__,_|\__|_|  \___/|_|  |_| |_| |_|\____|_|_|\___|_|\_\\___/| .__/ 
																	 |_|    
`

var iFunc functions.IFunctions
var iCPService cp_service.ICPService
var iClickupService clickup_service.IClickupService

func InitializeDependencyInjection() {
	iFunc = functions.GetFunctionsSingletonInstance()
	iCPService = cp_service.GetCPServiceSingletonInstance()
	iClickupService = clickup_service.GetClickupServiceSingletonInstance()
}

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

func MainAction(mainAction int) {

	iFunc.Log("...Starting ClickUp Automation...", true, variables_global.Config.ConfigGeneral.SaveLogInFile)

	for i := 0; i < len(variables_global.Config.Integrations); i++ {

		iFunc.Log("Found List "+variables_global.Config.Integrations[i].IntegrationName, true, variables_global.Config.ConfigGeneral.SaveLogInFile)
		iFunc.Log("Begin", true, variables_global.Config.ConfigGeneral.SaveLogInFile)

		if variables_global.Config.Integrations[i].OnlyCreateTask {
			iFunc.Log("List used only to create task", true, variables_global.Config.ConfigGeneral.SaveLogInFile)
			continue
		}

		list, error := iClickupService.ReturnList(variables_global.Config.Integrations[i].ClickUpListId)

		if error != nil {
			iFunc.Log("Error loading list "+variables_global.Config.Integrations[i].IntegrationName, true, variables_global.Config.ConfigGeneral.SaveLogInFile)
			continue
		}

		variables_global.Customer = variables_global.Config.Integrations[i]

		switch mainAction {
		case enum_main_action.TASKS_VERIFY:
			iClickupService.VerifyErrorsProjectWithStore(list)

		case enum_main_action.TASKS_UPDATE:
			iClickupService.UpdateProjectWithStore(list)
			iClickupService.UpdateTasksInDoneToClosedPSHierarchy(list, enum_clickup_type_ps_hierarchy.TASK)
			iClickupService.UpdateTasksInDoneToClosedPSHierarchy(list, enum_clickup_type_ps_hierarchy.STORE)
			iClickupService.UpdateTasksInDoneToClosedPSHierarchy(list, enum_clickup_type_ps_hierarchy.EPIC)

		case enum_main_action.TASKS_UPDATE_DONE_CLOSED:
			iClickupService.UpdateTasksInDoneToClosed(list)
		}

		iFunc.Log("Finish!", true, variables_global.Config.ConfigGeneral.SaveLogInFile)
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
			iCPService.InputSearchRequimentsPlatform()
		case 2:
			iCPService.InputSearchProjectTypesPlatform()
		// case 3:
		// 	service_conviso_platform.RetDeploys()
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
	label := iFunc.GetTextWithSpace("Label: ")

	//fmt.Print("Goal: ")
	goal := iFunc.GetTextWithSpace("Goal: ")

	//fmt.Print("Scope: ")
	scope := iFunc.GetTextWithSpace("Scope: ")

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

	createConvisoPlatform := type_platform.ProjectCreateRequestInput{variables_global.Customer.PlatformID,
		label, goal, scope, typeId,
		iFunc.ConvertStringToArrayInt(requirementIds),
		time.Now().Add(time.Hour * 24).Format("2006-01-02"), "1"}

	createProject, err := iCPService.AddPlatformProject(createConvisoPlatform)

	if err != nil || len(createProject.Data.CreateProject.Errors) > 0 {

		msgError := "Error CreateProject: "

		if err != nil {
			msgError = msgError + err.Error()
		}

		for i := 0; i < len(createProject.Data.CreateProject.Errors); i++ {
			msgError = msgError + " " + createProject.Data.CreateProject.Errors[i]
		}

		fmt.Println(msgError)
		return
	}

	customFieldUrlConvisoPlatform := type_clickup.CustomFieldRequest{
		variables_global.Config.ConfclickUp.CustomFieldPsCPLinkId,
		"https://app.convisoappsec.com/scopes/" + strconv.Itoa(variables_global.Customer.PlatformID) + "/projects/" + createProject.Data.CreateProject.Project.Id,
	}

	customFieldPSHierarchy := type_clickup.CustomFieldRequest{
		variables_global.Config.ConfclickUp.CustomFieldPsHierarchyId,
		strconv.Itoa(enum_clickup_type_ps_hierarchy.STORE),
	}

	customFieldPSTeam := type_clickup.CustomFieldRequest{
		variables_global.Config.ConfclickUp.CustomFieldPsTeamId,
		strconv.Itoa(enum_clickup_ps_team.CONSULTING),
	}

	customerOrder, err := iClickupService.RetClickUpDropDownPosition(variables_global.Customer.ClickUpListId, variables_global.Config.ConfclickUp.CustomFieldPsCustomerId,
		variables_global.Customer.ClickUpCustomerList)

	if err != nil {
		fmt.Println("Error customerOrder: Contact the system administrator")
		return
	}

	customFieldPSCustomer := type_clickup.CustomFieldRequest{
		variables_global.Config.ConfclickUp.CustomFieldPsCustomerId,
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
	taskMainClickup, err := iClickupService.TaskCreateRequest(
		type_clickup.TaskCreateRequest{
			createProject.Data.CreateProject.Project.Label,
			createProject.Data.CreateProject.Project.Scope,
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

		for i := 0; i < len(createProject.Data.CreateProject.Project.Activities); i++ {
			var convisoPlatformUrl bytes.Buffer
			convisoPlatformUrl.WriteString(variables_constant.CONVISO_PLATFORM_URL_BASE)
			convisoPlatformUrl.WriteString("scopes/")
			convisoPlatformUrl.WriteString(strconv.Itoa(variables_global.Customer.PlatformID))
			convisoPlatformUrl.WriteString("/projects/")
			convisoPlatformUrl.WriteString(createProject.Data.CreateProject.Project.Id)
			convisoPlatformUrl.WriteString("/project_requirements/")
			convisoPlatformUrl.WriteString(createProject.Data.CreateProject.Project.Activities[i].Id)

			customFieldUrlConvisoPlatformSubTask := type_clickup.CustomFieldRequest{
				variables_global.Config.ConfclickUp.CustomFieldPsCPLinkId,
				convisoPlatformUrl.String(),
			}

			customFieldTypeConsultingSubTask := type_clickup.CustomFieldRequest{
				variables_global.Config.ConfclickUp.CustomFieldPsHierarchyId,
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

			sanitizedHTMLTitle, err = html2text.FromString(createProject.Data.CreateProject.Project.Activities[i].Title)
			sanitizedHTMLDescription, err = html2text.FromString(createProject.Data.CreateProject.Project.Activities[i].Description)

			_, err := iClickupService.TaskCreateRequest(
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
				fmt.Println("Error CreateProject: Problem create ClickUp SubTask ", createProject.Data.CreateProject.Project.Activities[i].Title)
			}
		}
	}

	fmt.Println("Create Task Success!")
}

func InitialCheck() bool {
	ret := true

	err := error(nil)

	variables_global.Config, err = iFunc.LoadConfigsByYamlFile()

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

	if strings.EqualFold(variables_global.Config.ConfclickUp.CustomFieldPsCPLinkId, "") {
		variables_global.Config.ConfclickUp.CustomFieldPsCPLinkId = "ce46360c-373e-48a2-843f-eb3fd5bbf497"
	}

	if strings.EqualFold(variables_global.Config.ConfclickUp.CustomFieldPsHierarchyId, "") {
		variables_global.Config.ConfclickUp.CustomFieldPsHierarchyId = "addf1ada-84b4-494c-b640-45d0bb698181"
	}

	if strings.EqualFold(variables_global.Config.ConfclickUp.CustomFieldPsTeamId, "") {
		variables_global.Config.ConfclickUp.CustomFieldPsTeamId = "b30739d1-4169-41dc-b90b-e670a51b8545"
	}

	if strings.EqualFold(variables_global.Config.ConfclickUp.CustomFieldPsCustomerId, "") {
		variables_global.Config.ConfclickUp.CustomFieldPsCustomerId = "f3b4ecc4-737b-4040-a75d-664e89ad2f3a"
	}

	if strings.EqualFold(variables_global.Config.ConfclickUp.CustomFieldPsDeliveryPointsId, "") {
		variables_global.Config.ConfclickUp.CustomFieldPsDeliveryPointsId = "761d0b8d-5586-4c91-b861-1bc49210a0ee"
	}
}

func main() {
	/*
		TODO LIST
			separar as atualizações do cp e clickup, hoje tem uma variável, has update, mas deveria ter algo do tipo hasupdate cp e hasupcate clickup
			qdo não encontrar um cliente no campo PS Customer, não quebrar a aplicação, selecionar o primeiro da lista ou algo assim
			verificar possibilidade de melhorar a função de recuperar o customfield do clickup na função returntask
	*/

	InitializeDependencyInjection()

	//iFunc.WriteFile("teste", "teste tiago teste tiago")

	if !InitialCheck() {
		fmt.Println("You need to correct the above information before rerunning the application")
		fmt.Println("Press the Enter Key to finish!")
		fmt.Scanln()
		os.Exit(0)
	}

	SetDefaultValue()

	integrationJustVerify := flag.Bool("tv", false, "Verify if clickup tasks is ok")
	integrationUpdateTasks := flag.Bool("tu", false, "Update Conviso Platform and ClickUp Tasks")
	// integrationListTasksInProgress := flag.Bool("tsip", false, "List Clickup Stories In Progress")
	// integrationListTasksClosed := flag.Bool("tsd", false, "List Clickup Epics, Stories and Tasks in Closed")
	integrationUpdateTasksDone := flag.Bool("tud", false, "Change tasks done to closed")
	// deploy := flag.Bool("d", false, "See info about deploys")
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

	// if *integrationListTasksInProgress {
	// 	MainAction(enum_main_action.TASKS_INPROGRESS)
	// 	os.Exit(0)
	// }

	// if *integrationListTasksClosed {
	// 	MainAction(enum_main_action.TASKS_INCLOSED)
	// 	os.Exit(0)
	// }

	if *integrationUpdateTasksDone {
		MainAction(enum_main_action.TASKS_UPDATE_DONE_CLOSED)
		os.Exit(0)
	}

	// if *deploy {
	// 	service_conviso_platform.RetDeploys()
	// 	os.Exit(0)
	// }

	if *version {
		fmt.Println("Script Version: ", variables_constant.VERSION)
		os.Exit(0)
	}

	MainMenu()
}
