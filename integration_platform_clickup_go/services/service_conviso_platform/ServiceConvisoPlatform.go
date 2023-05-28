package ServiceConvisoPlatform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"golang.org/x/exp/slices"
	TypePlatform "integration.platform.clickup/types/type_platform"
	Functions "integration.platform.clickup/utils/functions"
	VariablesConstant "integration.platform.clickup/utils/variables_constant"
	VariablesGlobal "integration.platform.clickup/utils/variables_global"
)

const CONVISO_PLATFORM_PROJECT_TYPES = `
	query ProjectTypes($Page:Int, $ProjectType:String){
		projectTypes(page: $Page, limit: 10, params: {
		labelCont:$ProjectType
	}) {
		collection {
			code
			description
			id
			label
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

const CONVISO_PLATFORM_PROJECTS_QUERY = `
	query Projects($label: String){
		projects(page: 1, limit: 100, params: {
			labelEq: $label,
		}, sortBy: "id", descending: true) {
			collection {
				activities{
					id
					title,
					description
				}
				company{
					id
				}
				id
				label
				objective
				scope
				}
			metadata {
				currentPage
				totalPages
			}
		}
	}
`

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

func SearchRequimentsPlatform(reqSearch string) {
	var result TypePlatform.RequirementsResponse
	writer := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', tabwriter.AlignRight)
	fmt.Fprintf(writer, "\n %s\t%s\t", "Id", "Label")
	fmt.Fprintf(writer, "\n %s\t%s\t", "----", "----")

	var platformId int

	if VariablesGlobal.Customer.PlatformID == 0 {
		platformId = 11
	} else {
		platformId = VariablesGlobal.Customer.PlatformID
	}

	for i := 0; i <= result.Data.Playbooks.Metadata.TotalPages; i++ {
		parameters, _ := json.Marshal(TypePlatform.RequirementsParameters{CompanyId: platformId, Page: i + 1, Requirement: reqSearch})
		body, _ := json.Marshal(map[string]string{
			"query":     CONVISO_PLATFORM_REQUIREMENTS_QUERY,
			"variables": string(parameters),
		})

		payload := bytes.NewBuffer(body)
		req, err := http.NewRequest(http.MethodPost, VariablesConstant.CONVISO_PLATFORM_API_GRAPHQL, payload)
		if err != nil {
			fmt.Println("Error SearchRequimentsPlatform NewRequest: ", err.Error())
			return
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("x-api-key", os.Getenv("CONVISO_PLATFORM_TOKEN"))
		client := &http.Client{Timeout: time.Second * 10}
		resp, err := client.Do(req)
		defer req.Body.Close()
		if err != nil {
			fmt.Println("Error SearchRequimentsPlatform ClientDo: ", err.Error())
			return
		}
		data, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal([]byte(string(data)), &result)

		for i := 0; i < len(result.Data.Playbooks.Collection); i++ {
			fmt.Fprintf(writer, "\n %s\t%s\t",
				result.Data.Playbooks.Collection[i].Id,
				result.Data.Playbooks.Collection[i].Label)
		}
		writer.Flush()
		fmt.Println("")
		if result.Data.Playbooks.Metadata.CurrentPage != result.Data.Playbooks.Metadata.TotalPages {
			var input int
			fmt.Println("See next results? 0 - no; 1 - yes")
			fmt.Print("Enter the option: ")
			fmt.Scan(&input)
			if input == 0 {
				break
			}
		}
		result.Data.Playbooks.Metadata.TotalPages = result.Data.Playbooks.Metadata.TotalPages - 1
	}
}

func SearchProjectTypesPlatform(tpSearch string) {
	var result TypePlatform.ProjectTypesResponse
	writer := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', tabwriter.AlignRight)
	fmt.Fprintf(writer, "\n %s\t%s\t%s\t%s\t", "Id", "Label", "Code", "Description")
	fmt.Fprintf(writer, "\n %s\t%s\t%s\t%s\t", "----", "----", "----", "----")

	for i := 0; i <= result.Data.ProjectTypes.Metadata.TotalPages; i++ {
		parameters, _ := json.Marshal(TypePlatform.ProjectTypeParameters{Page: i + 1, ProjectType: tpSearch})
		body, _ := json.Marshal(map[string]string{
			"query":     CONVISO_PLATFORM_PROJECT_TYPES,
			"variables": string(parameters),
		})

		payload := bytes.NewBuffer(body)
		req, err := http.NewRequest(http.MethodPost, VariablesConstant.CONVISO_PLATFORM_API_GRAPHQL, payload)
		if err != nil {
			fmt.Println("Error SearchProjectTypesPlatform NewRequest: ", err.Error())
			return
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("x-api-key", os.Getenv("CONVISO_PLATFORM_TOKEN"))
		client := &http.Client{Timeout: time.Second * 10}
		resp, err := client.Do(req)
		defer req.Body.Close()
		if err != nil {
			fmt.Println("Error SearchProjectTypesPlatform ClientDo: ", err.Error())
			return
		}
		data, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal([]byte(string(data)), &result)
		for i := 0; i < len(result.Data.ProjectTypes.Collection); i++ {
			fmt.Fprintf(writer, "\n %s\t%s\t%s\t%s\t",
				result.Data.ProjectTypes.Collection[i].Id,
				result.Data.ProjectTypes.Collection[i].Label,
				result.Data.ProjectTypes.Collection[i].Code,
				result.Data.ProjectTypes.Collection[i].Description)
		}
		writer.Flush()
		fmt.Println("")
		if result.Data.ProjectTypes.Metadata.CurrentPage != result.Data.ProjectTypes.Metadata.TotalPages {
			var input int
			fmt.Println("See next results? 0 - no; 1 - yes")
			fmt.Print("Enter the option: ")
			fmt.Scan(&input)
			if input == 0 {
				break
			}
		}
		result.Data.ProjectTypes.Metadata.TotalPages = result.Data.ProjectTypes.Metadata.TotalPages - 1
	}
}

func InputSearchProjectTypesPlatform() {
	fmt.Print("Enter part of the project type: ")
	tpSearch := Functions.GetTextWithSpace()
	SearchProjectTypesPlatform(tpSearch)
}

func InputSearchRequimentsPlatform() {
	fmt.Print("Enter part of the requirement: ")
	reqSearch := Functions.GetTextWithSpace()
	SearchRequimentsPlatform(reqSearch)
}

func ConfirmProjectCreate(companyId int, label string) (TypePlatform.ProjectCollectionResponse, error) {
	var result TypePlatform.ProjectsResponse

	parameters, _ := json.Marshal(map[string]string{
		"label": label,
	})

	body, _ := json.Marshal(map[string]string{
		"query":     CONVISO_PLATFORM_PROJECTS_QUERY,
		"variables": string(parameters),
	})

	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, VariablesConstant.CONVISO_PLATFORM_API_GRAPHQL, payload)
	if err != nil {
		return TypePlatform.ProjectCollectionResponse{}, errors.New("Error ConfirmProjectCreate New Request " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", os.Getenv("CONVISO_PLATFORM_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		return TypePlatform.ProjectCollectionResponse{}, errors.New("Error ConfirmProjectCreate ClientDo " + err.Error())
	}
	data, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal([]byte(string(data)), &result)

	auxCompanyId, _ := strconv.Atoi(result.Data.Projects.Collection[0].Company.Id)

	if auxCompanyId != VariablesGlobal.Customer.PlatformID {
		return TypePlatform.ProjectCollectionResponse{}, errors.New("Different Company")
	}

	return result.Data.Projects.Collection[0], nil
}

func AddPlatformProject(inputParameters TypePlatform.ProjectCreateInputRequest) error {
	var tokenPlatform = os.Getenv("CONVISO_PLATFORM_TOKEN")

	projectParameters := TypePlatform.ProjectCreateRequest{inputParameters}

	parameters, _ := json.Marshal(projectParameters)
	body, _ := json.Marshal(map[string]string{
		"query":     CONVISO_PLATFORM_PROJECT_CREATE,
		"variables": string(parameters),
	})

	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, VariablesConstant.CONVISO_PLATFORM_API_GRAPHQL, payload)
	if err != nil {
		return errors.New("Error AddPlatformProject New Request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", tokenPlatform)
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		return errors.New("Error AddPlatformProject ClientDo: " + err.Error())
	}

	ioutil.ReadAll(resp.Body)

	return nil
}

func RequirementsId(text string) (int, error) {

	textSplit := strings.Split(text, "/")

	if !slices.Contains(textSplit, "project_requirements") {
		return 0, errors.New("Error RequirementsId Slice 0")
	}

	ret, error := strconv.Atoi(textSplit[len(textSplit)-1])

	if error != nil {
		return 0, errors.New("Error RequirementsId Not Integer")
	}

	return ret, nil
}

func ChangeActivitiesStatus(url string) {
	fmt.Println(url)

	RequirementsId(url)
}
