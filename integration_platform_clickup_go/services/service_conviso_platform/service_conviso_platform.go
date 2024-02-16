package service_conviso_platform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"integration_platform_clickup_go/types/type_clickup"
	"integration_platform_clickup_go/types/type_enum/enum_requirement_activity_status"
	"integration_platform_clickup_go/types/type_platform"
	"integration_platform_clickup_go/utils/functions"
	"integration_platform_clickup_go/utils/variables_constant"
	"integration_platform_clickup_go/utils/variables_global"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"golang.org/x/exp/slices"
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

const CONVISO_PLATFORM_DEPLOY = `
query DeploysByCompanyId($Page:Int){
	deploysByCompanyId(
	  id: "152"
	  initialDate: "2023-01-01T00:00:00-03:00"
	  finishDate: "2023-05-29T23:59:59-03:00"
	  page: $Page
	  limit: 1000
	) {
	  collection {
		changedApproximately
		changedLines
		createdAt
		currentCommit
		currentTag
		deployUrlCompare
		discardReason
		discardedId
		gauntletDiffUrl
		gauntletScanId
		gauntletSourceCodeId
		gitDiff
		id
		languages
		newLines
		previousCommit
		previousTag
		removedLines
		reviewed
		reviewedAt
		status
		totalProjectLines
		updatedAt,
		project{
		  label
		}
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

const CONVISO_PLATFORM_PROJECT_QUERY = `
	query Project($id: ID!)
	{
		project(id: $id) {
			activities{
				id
				title,
				description,
				status
			}
			company{
				id
			}
			id
			label
			objective
			scope,
			status
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

const CONVISO_PLATFORM_UPDATE_REQUIREMENTS_ACTIVITY_START = `
	mutation  UpdateActivityStatusToStart($input:UpdateActivityStatusToStartInput!)
	{
		updateActivityStatusToStart(input: $input) 
		{
			clientMutationId
			errors
		}
	}
`

const CONVISO_PLATFORM_UPDATE_REQUIREMENTS_ACTIVITY_RESTART = `
	mutation  UpdateActivityStatusToRestart($input:UpdateActivityStatusToRestartInput!)
	{
		updateActivityStatusToRestart(input: $input) 
		{
			clientMutationId
			errors
		}
	}
`

const CONVISO_PLATFORM_UPDATE_REQUIREMENTS_ACTIVITY_FINISH = `
	mutation  UpdateActivityStatusToFinish($input:UpdateActivityStatusToFinishInput!)
	{
		updateActivityStatusToFinish(input: $input) 
		{
			clientMutationId
			errors
		}
	}
`

func RetDeploys() {
	var result type_platform.DeployTypeResponse

	reviewNewLine := 0
	reviewRemovedLine := 0
	reviewChangedLine := 0
	numDeploysReviewed := 0
	numDeploysNotReviewed := 0

	newLine := 0
	removedLine := 0
	changedLine := 0

	for i := 0; i <= result.Data.DeployTypeData.Metadata.TotalPages; i++ {
		parameters, _ := json.Marshal(type_platform.PageParameters{Page: i + 1})
		body, _ := json.Marshal(map[string]string{
			"query":     CONVISO_PLATFORM_DEPLOY,
			"variables": string(parameters),
		})

		payload := bytes.NewBuffer(body)
		req, err := http.NewRequest(http.MethodPost, variables_constant.CONVISO_PLATFORM_API_GRAPHQL, payload)
		if err != nil {
			fmt.Println("Error RetDeploys NewRequest: ", err.Error())
			return
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("x-api-key", os.Getenv("CONVISO_PLATFORM_TOKEN"))
		client := &http.Client{Timeout: time.Second * 30}
		resp, err := client.Do(req)
		defer req.Body.Close()
		if err != nil {
			fmt.Println("Error RetDeploys ClientDo: ", err.Error())
			return
		}
		data, _ := io.ReadAll(resp.Body)

		json.Unmarshal([]byte(string(data)), &result)
		for i := 0; i < len(result.Data.DeployTypeData.Collection); i++ {

			if result.Data.DeployTypeData.Collection[i].Reviewed {
				reviewChangedLine = reviewChangedLine + result.Data.DeployTypeData.Collection[i].ChangedLines
				reviewRemovedLine = reviewRemovedLine + result.Data.DeployTypeData.Collection[i].RemovedLines
				reviewNewLine = reviewNewLine + result.Data.DeployTypeData.Collection[i].NewLines
				numDeploysReviewed = numDeploysReviewed + 1
			} else {
				changedLine = changedLine + result.Data.DeployTypeData.Collection[i].ChangedLines
				removedLine = removedLine + result.Data.DeployTypeData.Collection[i].RemovedLines
				newLine = newLine + result.Data.DeployTypeData.Collection[i].NewLines
				numDeploysNotReviewed = numDeploysNotReviewed + 1

			}
		}
		result.Data.DeployTypeData.Metadata.TotalPages = result.Data.DeployTypeData.Metadata.TotalPages - 1

		fmt.Println(strconv.Itoa(i), "/", strconv.Itoa(result.Data.DeployTypeData.Metadata.TotalPages))
	}

	println("reviewChangedLine = " + strconv.Itoa(reviewChangedLine))
	println("reviewRemovedLine = " + strconv.Itoa(reviewRemovedLine))
	println("reviewNewLine = " + strconv.Itoa(reviewNewLine))
	println("changedLine = " + strconv.Itoa(changedLine))
	println("removedLine = " + strconv.Itoa(removedLine))
	println("newLine = " + strconv.Itoa(newLine))
	println("Total Deploys = " + strconv.Itoa(numDeploysReviewed+numDeploysNotReviewed))
	println("Deploys not Reviewed = " + strconv.Itoa(numDeploysNotReviewed))
	println("Deploys Reviewed = " + strconv.Itoa(numDeploysReviewed))
}

func SearchRequimentsPlatform(reqSearch string) {
	var result type_platform.RequirementsResponse
	writer := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', tabwriter.AlignRight)
	fmt.Fprintf(writer, "\n %s\t%s\t", "Id", "Label")
	fmt.Fprintf(writer, "\n %s\t%s\t", "----", "----")

	var platformId int

	if variables_global.Customer.PlatformID == 0 {
		platformId = 11
	} else {
		platformId = variables_global.Customer.PlatformID
	}

	for i := 0; i <= result.Data.Playbooks.Metadata.TotalPages; i++ {
		parameters, _ := json.Marshal(type_platform.RequirementsParameters{CompanyId: platformId, Page: i + 1, Requirement: reqSearch})
		body, _ := json.Marshal(map[string]string{
			"query":     CONVISO_PLATFORM_REQUIREMENTS_QUERY,
			"variables": string(parameters),
		})

		payload := bytes.NewBuffer(body)
		req, err := http.NewRequest(http.MethodPost, variables_constant.CONVISO_PLATFORM_API_GRAPHQL, payload)
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

		data, _ := io.ReadAll(resp.Body)

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
	var result type_platform.ProjectTypesResponse
	writer := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', tabwriter.AlignRight)
	fmt.Fprintf(writer, "\n %s\t%s\t%s\t%s\t", "Id", "Label", "Code", "Description")
	fmt.Fprintf(writer, "\n %s\t%s\t%s\t%s\t", "----", "----", "----", "----")

	for i := 0; i <= result.Data.ProjectTypes.Metadata.TotalPages; i++ {
		parameters, _ := json.Marshal(type_platform.ProjectTypeParameters{Page: i + 1, ProjectType: tpSearch})
		body, _ := json.Marshal(map[string]string{
			"query":     CONVISO_PLATFORM_PROJECT_TYPES,
			"variables": string(parameters),
		})

		payload := bytes.NewBuffer(body)
		req, err := http.NewRequest(http.MethodPost, variables_constant.CONVISO_PLATFORM_API_GRAPHQL, payload)
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
		data, _ := io.ReadAll(resp.Body)

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
	tpSearch := functions.GetTextWithSpace("Enter part of the project type: ")
	SearchProjectTypesPlatform(tpSearch)
}

func InputSearchRequimentsPlatform() {
	reqSearch := functions.GetTextWithSpace("Enter part of the requirement: ")
	SearchRequimentsPlatform(reqSearch)
}

func ConfirmProjectCreate(companyId int, label string) (type_platform.ProjectCollectionResponse, error) {
	var result type_platform.ProjectsResponse

	parameters, _ := json.Marshal(map[string]string{
		"label": label,
	})

	body, _ := json.Marshal(map[string]string{
		"query":     CONVISO_PLATFORM_PROJECTS_QUERY,
		"variables": string(parameters),
	})

	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, variables_constant.CONVISO_PLATFORM_API_GRAPHQL, payload)
	if err != nil {
		return type_platform.ProjectCollectionResponse{}, errors.New("Error ConfirmProjectCreate New Request " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", os.Getenv("CONVISO_PLATFORM_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		return type_platform.ProjectCollectionResponse{}, errors.New("Error ConfirmProjectCreate ClientDo " + err.Error())
	}
	data, _ := io.ReadAll(resp.Body)

	json.Unmarshal([]byte(string(data)), &result)

	auxCompanyId, _ := strconv.Atoi(result.Data.Projects.Collection[0].Company.Id)

	if auxCompanyId != variables_global.Customer.PlatformID {
		return type_platform.ProjectCollectionResponse{}, errors.New("Different Company")
	}

	return result.Data.Projects.Collection[0], nil
}

func GetProject(id int) (type_platform.Project, error) {
	var result type_platform.ProjectResponse

	parameters, _ := json.Marshal(map[string]int{
		"id": id,
	})

	body, _ := json.Marshal(map[string]string{
		"query":     CONVISO_PLATFORM_PROJECT_QUERY,
		"variables": string(parameters),
	})

	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, variables_constant.CONVISO_PLATFORM_API_GRAPHQL, payload)
	if err != nil {
		return type_platform.Project{}, errors.New("Error GetProject New Request " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", os.Getenv("CONVISO_PLATFORM_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		return type_platform.Project{}, errors.New("Error GetProject ClientDo " + err.Error())
	}
	data, _ := io.ReadAll(resp.Body)

	json.Unmarshal([]byte(string(data)), &result)

	return result.Data.Project, nil
}

func AddPlatformProject(inputParameters type_platform.ProjectCreateInputRequest) error {
	var tokenPlatform = os.Getenv("CONVISO_PLATFORM_TOKEN")

	projectParameters := type_platform.ProjectCreateRequest{inputParameters}

	parameters, _ := json.Marshal(projectParameters)
	body, _ := json.Marshal(map[string]string{
		"query":     CONVISO_PLATFORM_PROJECT_CREATE,
		"variables": string(parameters),
	})

	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, variables_constant.CONVISO_PLATFORM_API_GRAPHQL, payload)
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

	io.ReadAll(resp.Body)

	return nil
}

func RetQueryUpdateRequirementsActivityStatus(action int) string {
	switch action {
	case enum_requirement_activity_status.START:
		return CONVISO_PLATFORM_UPDATE_REQUIREMENTS_ACTIVITY_START
	case enum_requirement_activity_status.FINISH:
		return CONVISO_PLATFORM_UPDATE_REQUIREMENTS_ACTIVITY_FINISH
	case enum_requirement_activity_status.RESTART:
		return CONVISO_PLATFORM_UPDATE_REQUIREMENTS_ACTIVITY_RESTART
	default:
		return ""
	}
}

func ChangeActivitiesStatusGraphQl(activityId int, action int) error {
	var tokenPlatform = os.Getenv("CONVISO_PLATFORM_TOKEN")

	input := type_platform.UpdateRequirementsActivityStatusInputRequest{activityId}

	activiTyStatusParameters := type_platform.UpdateRequirementsActivityStatusRequest{input}

	parameters, _ := json.Marshal(activiTyStatusParameters)
	body, _ := json.Marshal(map[string]string{
		"query":     RetQueryUpdateRequirementsActivityStatus(action),
		"variables": string(parameters),
	})

	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, variables_constant.CONVISO_PLATFORM_API_GRAPHQL, payload)
	if err != nil {
		return errors.New("Error ChangeActivitiesStatusGraphQl New Request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", tokenPlatform)
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		return errors.New("Error ChangeActivitiesStatusGraphQl ClientDo: " + err.Error())
	}

	io.ReadAll(resp.Body)

	return nil
}

func RetActivityIdCustomField(text string) (int, error) {

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

func RetProjectIdCustomField(text string) (int, error) {

	textSplit := strings.Split(text, "/")

	if !slices.Contains(textSplit, "projects") {
		return 0, errors.New("Error RetProjectIdCustomField don't have projects")
	}

	ret, error := strconv.Atoi(textSplit[len(textSplit)-1])

	if error != nil {
		return 0, errors.New("Error RetProjectIdCustomField Not Integer")
	}

	return ret, nil
}

func UpdateActivityRequirement(task type_clickup.TaskResponse, project type_platform.Project) error {
	if project.Id != "" {
		activityId, error := RetActivityIdCustomField(task.CustomField.PSConvisoPlatformLink)

		if error == nil {
			idxActivity := slices.IndexFunc(project.Activities, func(a type_platform.ActivityCollectionResponse) bool { return a.Id == strconv.Itoa(activityId) })
			if idxActivity != -1 {
				switch strings.ToLower(task.Status.Status) {
				case "backlog", "to do":
					if strings.ToLower(project.Activities[idxActivity].Status) == "done" {
						ChangeActivitiesStatusGraphQl(activityId, enum_requirement_activity_status.RESTART)
					}
				case "in progress":
					if strings.ToLower(project.Activities[idxActivity].Status) == "not_started" {
						ChangeActivitiesStatusGraphQl(activityId, enum_requirement_activity_status.START)
					}
				case "done":
					if strings.ToLower(project.Activities[idxActivity].Status) == "not_started" {
						ChangeActivitiesStatusGraphQl(activityId, enum_requirement_activity_status.START)
						ChangeActivitiesStatusGraphQl(activityId, enum_requirement_activity_status.FINISH)
					} else if strings.ToLower(project.Activities[idxActivity].Status) == "in_progress" {
						ChangeActivitiesStatusGraphQl(activityId, enum_requirement_activity_status.FINISH)
					}
				}
			}
		} else {
			return error
		}
	}

	return nil
}

func UpdateProjectRest(request type_clickup.TaskRequestStore, cpProjectId string, timeEstimate int64) error {

	data := url.Values{}
	data.Set("project_status_id", RetNewStatus(request.Status)) //1 - running; 3 done; 4-planned
	data.Set("accept", "1")
	data.Set("project_history_started_at_aux", (time.UnixMilli(request.StartDate)).Format("02/01/2006"))
	data.Set("project_history_estimated_hours", strconv.Itoa(int(timeEstimate)))
	data.Set("project_history_delivery_date", (time.UnixMilli(request.DueDate)).Format("2006-01-02"))
	data.Set("project_history_delivery_date_aux", (time.UnixMilli(request.DueDate)).Format("02/01/2006"))

	var tokenPlatform = os.Getenv(variables_constant.CONVISO_PLATFORM_TOKEN_NAME)

	var convisoPlatformUrl bytes.Buffer
	convisoPlatformUrl.WriteString(variables_constant.CONVISO_PLATFORM_URL_BASE)
	convisoPlatformUrl.WriteString("scopes/")
	convisoPlatformUrl.WriteString(strconv.Itoa(variables_global.Customer.PlatformID))
	convisoPlatformUrl.WriteString("/projects/")
	convisoPlatformUrl.WriteString(cpProjectId)
	convisoPlatformUrl.WriteString("/update_status")

	req, err := http.NewRequest(http.MethodPost, convisoPlatformUrl.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return errors.New("Error UpdateProjectRest New Request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Set("X-Armature-Api-Key", tokenPlatform)
	req.Header.Set("Accept", "*/*")

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		return errors.New("Error UpdateProjectRest ClientDo: " + err.Error() + "StatusCode: " + resp.Status)
	}

	return nil
}

func RetNewStatus(statusTask string) string {
	//1 - running; 3 done; 4-planned
	newReturn := statusTask
	switch statusTask {
	case "backlog", "to do":
		newReturn = "4"
	case "in progress":
		newReturn = "1"
	case "done":
		newReturn = "3"
	}
	return newReturn
}
