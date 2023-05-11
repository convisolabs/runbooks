package main

import (
	"fmt"
	"time"

	ServicesClickup "integration.platform.clickup/services/service_clickup"
	ServiceConvisoPlatform "integration.platform.clickup/services/service_conviso_platform"
	TypeClickup "integration.platform.clickup/types/type_clickup"
	TypePlatform "integration.platform.clickup/types/type_platform"
	Functions "integration.platform.clickup/utils/functions"
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

// func MenuRequirementsSearch() {
// 	var input int
// 	for ok := true; ok; ok = (input != 0) {
// 		fmt.Println("-----Menu Requirements Search-----")
// 		fmt.Println("Project Selected: ", VariablesGlobal.Customer.Name)
// 		fmt.Println("0 - Previous Menu")
// 		fmt.Println("1 - Search Requirements")
// 		fmt.Print("Enter the option: ")
// 		n, err := fmt.Scan(&input)
// 		if n < 1 || err != nil {
// 			fmt.Println("Invalid Input")
// 			break
// 		}
// 		switch input {
// 		case 0:
// 			break
// 		case 1:
// 			ServiceConvisoPlatform.InputSearchRequimentsPlatform()
// 		default:
// 			fmt.Println("Invalid Input")
// 		}
// 	}
// }

func MenuClickup() {
	var input int
	for ok := true; ok; ok = (input != 0) {
		fmt.Println("-----Menu Clickup-----")
		fmt.Println("0 - Previous Menu")
		fmt.Println("1 - Verification Tasks Clickup")
		fmt.Println("2 - Update Tasks Clickup")
		fmt.Print("Enter the option: ")
		n, err := fmt.Scan(&input)
		if n < 1 || err != nil {
			fmt.Println("Invalid Input")
		}
		switch input {
		case 0:
			break
		case 1:
			ServicesClickup.ClickUpAutomation(true)
		case 2:
			ServicesClickup.ClickUpAutomation(false)
		default:
			fmt.Println("Invalid Input")
		}
	}
}

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
		fmt.Println("Project Selected: ", VariablesGlobal.Customer.Name)
		fmt.Println("0 - Exit")
		fmt.Println("1 - Menu Clickup")
		fmt.Println("2 - Menu Setup")
		fmt.Println("3 - Create Project Conviso Platform/ClickUp")
		fmt.Println("4 - Menu Search Conviso Platform")

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
			MenuClickup()
		case 2:
			MenuSetupConfig()
		case 3:
			if VariablesGlobal.Customer.PlatformID == 0 {
				fmt.Println("No Project Selected!")
			} else {
				CreateProject()
			}
		case 4:
			MenuSearchConvisoPlatform()
		default:
			fmt.Println("Invalid Input")
		}
	}
}

func CreateProject() {

	// taskMainClickup, err := ServicesClickup.TaskCreateRequest(
	// 	TypeClickup.TaskCreateRequest{
	// 		"project.Label",
	// 		"project.Scope",
	// 		"backlog",
	// 		true,
	// 		""})

	playbookIds := ""
	typeId := 10

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

	// //primeiro lugar a criar no conviso platforme e depois no clickup
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
	}

	customFields := []TypeClickup.CustomFieldRequest{
		TypeClickup.CustomFieldRequest{
			"8e2863f4-e11f-409c-a373-893bc12200fb",
			"https://app.convisoappsec.com/scopes/" + string(VariablesGlobal.Customer.PlatformID) + "/projects/" + project.Id,
		},
		TypeClickup.CustomFieldRequest{
			"664816bc-a899-45ec-9801-5a1e5be9c5f6",
			"0",
		},
		TypeClickup.CustomFieldRequest{
			"4493a404-3ef7-4d7a-91e4-830ebc666353",
			"1",
		},
	}

	//create main
	taskMainClickup, err := ServicesClickup.TaskCreateRequest(
		TypeClickup.TaskCreateRequest{
			project.Label,
			project.Scope,
			"backlog",
			true,
			"",
			"",
			customFields})

	if err != nil {
		fmt.Println("problem")
	}

	for i := 0; i < len(project.Activities); i++ {
		_, err := ServicesClickup.TaskCreateRequest(
			TypeClickup.TaskCreateRequest{
				project.Activities[i].Title,
				project.Activities[i].Description,
				"backlog",
				true,
				taskMainClickup.Id,
				taskMainClickup.Id,
				customFields})
		if err != nil {
			fmt.Println("problem for ", i)
		}
	}

	//create subtasks
	fmt.Println(project)
	fmt.Println(taskMainClickup)

}

func main() {
	//próximas tarefas
	// listar os tipos do conviso platform
	// setar os campos customizáveis na função create
	// parametrizar tudo para utilização da função flags do golang

	// 	var searchRequiments string

	// 	flag.StringVar(&searchRequiments, "sr", "", "Search Conviso Platform Requirements")

	// 	flag.Parse()

	// 	if searchRequiments != "" {
	// 		ServiceConvisoPlatform.SearchRequimentsPlatform(searchRequiments)
	// 	} else {
	MainMenu()
	// }
}
