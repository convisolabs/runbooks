package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/exp/slices"

	ServicesClickup "integration.platform.clickup/services/service_clickup"
	ServiceConvisoPlatform "integration.platform.clickup/services/service_conviso_platform"
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
			ServiceConvisoPlatform.SearchRequimentsPlatform()
		default:
			fmt.Println("Invalid Input")
		}
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
			fmt.Println("3 - Menu Setup")
			fmt.Println("4 - Menu Search Requirements Conviso Platform")
			fmt.Println("5 - Create Project Conviso Platform/ClickUp")

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
			case 3:
				MenuSetupConfig()
			case 4:
				MenuRequirementsSearch()
			case 5:
				if VariablesGlobal.Customer.PlatformID == 0 {
					fmt.Println("Nenhum projeto selecionado!")
				} else {
					CreateProject()
				}

			default:
				fmt.Println("Invalid Input")
			}
		}
	}
}

func CreateProject() {

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

	fmt.Println(project)

	//precisa pegar o retorno do conviso platform e criar no clickup

}

func main() {
	MainMenu()
	//AddPlatformProject()
}
