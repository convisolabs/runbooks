package clickup_service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	cp_service "integration_platform_clickup_go/services/cp"
	"integration_platform_clickup_go/types/type_clickup"
	"integration_platform_clickup_go/types/type_enum/enum_clickup_ps_team"
	"integration_platform_clickup_go/types/type_enum/enum_clickup_type_ps_hierarchy"
	"integration_platform_clickup_go/types/type_platform"
	"integration_platform_clickup_go/utils/functions"
	"integration_platform_clickup_go/utils/variables_constant"
	"integration_platform_clickup_go/utils/variables_global"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

type ClickupService struct {
	HttpHeaders map[string]string
	Functions   functions.IFunctions
	CPService   cp_service.ICPService
}

func ClickupServiceNew(HttpHeaders map[string]string, Functions functions.IFunctions, CPService cp_service.ICPService) IClickupService {
	return &ClickupService{
		HttpHeaders: HttpHeaders,
		Functions:   Functions,
		CPService:   CPService,
	}
}

// var globalClickupHeaders map[string]string

// func init() {
// 	globalClickupHeaders = map[string]string{
// 		"Authorization": os.Getenv(variables_constant.CLICKUP_TOKEN_NAME),
// 	}
// }

func (f *ClickupService) RetAssigness(assignees []type_clickup.AssigneeField) string {
	ret := ""

	for i := 0; i < len(assignees); i++ {
		ret = ret + assignees[i].Username + ";"
	}

	return ret
}

func (f *ClickupService) RetClickUpDropDownPosition(clickupListId string, clickupFieldId string, searchValue string) (int, error) {
	result := -1
	customFieldsResponse, err := f.RetAllCustomFieldByList(clickupListId)

	if err != nil {
		return result, errors.New("Error RetClickUpDropDownPosition RequestCustomField: " + err.Error())
	}

	for i := 0; i < len(customFieldsResponse.Fields); i++ {

		if customFieldsResponse.Fields[i].Id == clickupFieldId {
			for j := 0; j < len(customFieldsResponse.Fields[i].TypeConfig.Options); j++ {
				if strings.EqualFold(customFieldsResponse.Fields[i].TypeConfig.Options[j].Name, searchValue) {
					return customFieldsResponse.Fields[i].TypeConfig.Options[j].OrderIndex, nil
				}
			}
		}
	}
	return result, nil
}

// func RetClickUpDropDownOptionName(clickupFieldId string, order int) (string, error) {
// 	result := ""
// 	customFieldsResponse, err := RetAllCustomFieldByList(clickupFieldId)

// 	if err != nil {
// 		return result, errors.New("Error RetClickUpDropDownOptionName RequestCustomField: " + err.Error())
// 	}

// 	for i := 0; i < len(customFieldsResponse.Fields); i++ {

// 		if customFieldsResponse.Fields[i].Id == clickupFieldId {
// 			for j := 0; j < len(customFieldsResponse.Fields[i].TypeConfig.Options); j++ {
// 				if customFieldsResponse.Fields[i].TypeConfig.Options[j].OrderIndex == order {
// 					return customFieldsResponse.Fields[i].TypeConfig.Options[j].Name, nil
// 				}
// 			}
// 		}
// 	}
// 	return result, nil
// }

func (f *ClickupService) RetCustomFieldValueString(customFieldId string, customFields []type_clickup.CustomField) string {
	for i := 0; i < len(customFields); i++ {
		if customFields[i].Id == customFieldId {
			if customFields[i].Value == nil {
				return ""
			} else {
				return customFields[i].Value.(string)
			}
		}
	}
	return ""
}

func (f *ClickupService) RetCustomFieldPSTeam(customFields []type_clickup.CustomField) string {
	for i := 0; i < len(customFields); i++ {
		if strings.EqualFold(customFields[i].Id, variables_global.Config.ConfclickUp.CustomFieldPsTeamId) {
			if customFields[i].Value == nil {
				return ""
			} else {
				return enum_clickup_ps_team.ToString(int(customFields[i].Value.(float64)))
			}
		}
	}
	return ""
}

func (f *ClickupService) RetCustomFieldPSCustomer(customFields []type_clickup.CustomField) string {
	for i := 0; i < len(customFields); i++ {
		if strings.EqualFold(customFields[i].Id, variables_global.Config.ConfclickUp.CustomFieldPsCustomerId) {
			if customFields[i].Value == nil {
				return ""
			} else {
				if len(customFields[i].TypeConfig.Options) > int(customFields[i].Value.(float64)) {
					return customFields[i].TypeConfig.Options[int(customFields[i].Value.(float64))].Name
				} else {
					return ""
				}
			}
		}
	}
	return ""
}

func (f *ClickupService) RetCustomFieldTypeConsulting(customFields []type_clickup.CustomField) int {
	for i := 0; i < len(customFields); i++ {
		if strings.EqualFold(customFields[i].Id, variables_global.Config.ConfclickUp.CustomFieldPsCustomerId) {
			if customFields[i].Value == nil {
				return -1
			} else {
				return int(customFields[i].Value.(float64))
			}
		}
	}
	return 0
}

// // func RetCustomFieldTeam(customFields []type_clickup.CustomField) []string {
// // 	for i := 0; i < len(customFields); i++ {
// // 		if customFields[i].Id == variables_constant.CLICKUP_CUSTOM_FIELD_PS_TEAM_ID {
// // 			if customFields[i].Value == nil {
// // 				return []string{}
// // 			} else {
// // 				aInterface := customFields[i].Value.([]interface{})
// // 				aString := make([]string, len(aInterface))
// // 				for i, v := range aInterface {
// // 					aString[i] = v.(string)
// // 				}
// // 				return aString
// // 			}
// // 		}
// // 	}
// // 	return []string{}
// // }

func (f *ClickupService) RetAllCustomFieldByList(listId string) (type_clickup.CustomFieldsResponse, error) {
	var result type_clickup.CustomFieldsResponse
	var urlGetTasks bytes.Buffer
	urlGetTasks.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlGetTasks.WriteString("list/")
	urlGetTasks.WriteString(listId)
	urlGetTasks.WriteString("/field")

	response, err := f.Functions.HttpRequestRetry(http.MethodGet, urlGetTasks.String(), f.HttpHeaders, nil, *variables_global.Config.ConfclickUp.HttpAttempt)

	if err != nil {
		return result, errors.New("Error RetCustomerPosition: " + err.Error())
	}

	data, _ := io.ReadAll(response.Body)

	json.Unmarshal([]byte(string(data)), &result)

	return result, nil
}

// func VerifyTasks(taskEpic type_clickup.TaskResponse) error {

// 	if len(taskEpic.LinkedTasks) == 0 {
// 		fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", "Nenhuma subtarefa lincada")
// 	}

// 	for k := 0; k < len(taskEpic.LinkedTasks); k++ {
// 		auxTaskId := ""

// 		if taskEpic.Id == taskEpic.LinkedTasks[k].LinkId {
// 			auxTaskId = taskEpic.LinkedTasks[k].TaskId
// 		} else {
// 			auxTaskId = taskEpic.LinkedTasks[k].LinkId
// 		}

// 		taskAux, err := ReturnTask(auxTaskId)
// 		if err != nil {
// 			return errors.New("Error taskAux: " + err.Error())
// 		}

// 		if strings.ToLower(taskAux.Status.Status) != "backlog" && strings.ToLower(taskAux.Status.Status) != "canceled" && strings.ToLower(taskAux.Status.Status) != "blocked" {
// 			if taskAux.DueDate == "" {
// 				fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", taskAux.Name, " :: ", "DueDate empty", " :: ", taskAux.Url)
// 			}

// 			if taskAux.StartDate == "" {
// 				fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", taskAux.Name, " :: ", "StartDate empty", " :: ", taskAux.Url)
// 			}

// 			if taskAux.TimeEstimate == 0 {
// 				fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", taskAux.Name, " :: ", "TimeEstimate empty", " :: ", taskAux.Url)
// 			}

// 			if taskAux.Status.Status == "done" && taskAux.TimeSpent == 0 {
// 				fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", taskAux.Name, " :: ", "TimeSpent empty", " :: ", taskAux.Url)
// 			}
// 		}
// 	}

// 	return nil
// }

func (f *ClickupService) ReturnTasks(listId string, searchTasks type_clickup.SearchTask) (type_clickup.TasksResponse, error) {
	var resultTasks type_clickup.TasksResponse
	var urlGetTasks bytes.Buffer

	urlGetTasks.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlGetTasks.WriteString("list/")
	urlGetTasks.WriteString(listId)
	urlGetTasks.WriteString("/task?custom_fields=[")
	urlGetTasks.WriteString("{\"field_id\":\"")
	urlGetTasks.WriteString(variables_global.Config.ConfclickUp.CustomFieldPsHierarchyId)
	urlGetTasks.WriteString("\",\"operator\":\"=\",\"value\":\"")
	urlGetTasks.WriteString(strconv.FormatInt(int64(searchTasks.TaskType), 10))
	urlGetTasks.WriteString("\"}")
	urlGetTasks.WriteString("]")

	if searchTasks.IncludeClosed {
		urlGetTasks.WriteString("&include_closed=true")
	}

	if searchTasks.DateUpdatedGt > 0 {
		urlGetTasks.WriteString("&date_updated_gt=")
		urlGetTasks.WriteString(strconv.FormatInt(searchTasks.DateUpdatedGt, 10))
	}

	if searchTasks.SubTasks {
		urlGetTasks.WriteString("&subtasks=true")
	}

	if !strings.EqualFold(searchTasks.TaskStatuses, "") {
		urlGetTasks.WriteString("&statuses[]=")
		urlGetTasks.WriteString(searchTasks.TaskStatuses)
	}

	urlGetTasks.WriteString("&page=")
	urlGetTasks.WriteString(strconv.FormatInt(int64(searchTasks.Page), 10))

	response, err := f.Functions.HttpRequestRetry(http.MethodGet, urlGetTasks.String(), f.HttpHeaders, nil, *variables_global.Config.ConfclickUp.HttpAttempt)

	if err != nil {
		return resultTasks, errors.New("Error ReturnTasks: " + err.Error())
	}

	data, _ := io.ReadAll(response.Body)

	json.Unmarshal([]byte(string(data)), &resultTasks)

	return resultTasks, nil
}

func (f *ClickupService) ReturnTask(taskId string) (type_clickup.TaskResponse, error) {
	var task type_clickup.TaskResponse
	var urlGetTask bytes.Buffer
	urlGetTask.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlGetTask.WriteString("task/")
	urlGetTask.WriteString(taskId)
	urlGetTask.WriteString("?include_subtasks=true")

	response, err := f.Functions.HttpRequestRetry(http.MethodGet, urlGetTask.String(), f.HttpHeaders, nil, *variables_global.Config.ConfclickUp.HttpAttempt)
	if err != nil {
		return task, errors.New("Error ReturnTask: " + err.Error())
	}

	data, _ := io.ReadAll(response.Body)

	json.Unmarshal([]byte(string(data)), &task)

	//add customFields
	task.CustomField.PSProjectHierarchy = f.RetCustomFieldTypeConsulting(task.CustomFields)
	task.CustomField.PSConvisoPlatformLink = f.RetCustomFieldValueString(variables_global.Config.ConfclickUp.CustomFieldPsCPLinkId, task.CustomFields)
	task.CustomField.PSTeam = f.RetCustomFieldPSTeam(task.CustomFields)
	task.CustomField.PSCustomer = f.RetCustomFieldPSCustomer(task.CustomFields)
	task.CustomField.PSDeliveryPoints = f.RetCustomFieldValueString(variables_global.Config.ConfclickUp.CustomFieldPsDeliveryPointsId, task.CustomFields)

	return task, nil
}

func (f *ClickupService) ReturnList(listId string) (type_clickup.ListResponse, error) {
	var list type_clickup.ListResponse
	var urlGetList bytes.Buffer
	urlGetList.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlGetList.WriteString("list/")
	urlGetList.WriteString(listId)
	urlGetList.WriteString("?include_subtasks=true")

	response, err := f.Functions.HttpRequestRetry(http.MethodGet, urlGetList.String(), f.HttpHeaders, nil, *variables_global.Config.ConfclickUp.HttpAttempt)

	if err != nil {
		return list, errors.New("Error ReturnList: " + err.Error())
	}

	data, _ := io.ReadAll(response.Body)

	json.Unmarshal([]byte(string(data)), &list)

	return list, nil
}

func (f *ClickupService) RetNewStatus(statusTask string, statusSubTask string) (string, bool) {

	newReturn := statusTask
	hasUpdate := false

	switch strings.ToLower(statusSubTask) {
	case "backlog":
		break
	case "to do":
		if statusTask == "backlog" {
			newReturn = "to do"
			hasUpdate = true
		}
	case "in progress", "done", "closed":
		if statusTask == "backlog" || statusTask == "to do" || statusTask == "blocked" {
			newReturn = "in progress"
			hasUpdate = true
		}
	}
	return newReturn, hasUpdate
}

// func RequestPutTask(taskId string, request type_clickup.TaskRequest) error {

// 	var urlPutTask bytes.Buffer
// 	urlPutTask.WriteString(variables_constant.CLICKUP_API_URL_BASE)
// 	urlPutTask.WriteString("task/")
// 	urlPutTask.WriteString(taskId)

// 	body, _ := json.Marshal(request)

// 	payload := bytes.NewBuffer(body)

// 	time.Sleep(time.Second)

// 	req, err := http.NewRequest(http.MethodPut, urlPutTask.String(), payload)
// 	if err != nil {
// 		return errors.New("Error RequestPutTask request: " + err.Error())
// 	}

// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
// 	client := &http.Client{Timeout: time.Second * 10}
// 	resp, err := client.Do(req)
// 	defer req.Body.Close()

// 	if resp.StatusCode != 200 {
// 		return errors.New("Error RequestPutTask StatusCode: " + http.StatusText(resp.StatusCode))
// 	}

// 	if err != nil {
// 		return errors.New("Error RequestPutTask response: " + err.Error())
// 	}

// 	io.ReadAll(resp.Body)

// 	return nil
// }

func (f *ClickupService) RequestPutTaskStatus(taskId string, request type_clickup.TaskRequestStatus) error {

	var urlPutTask bytes.Buffer
	urlPutTask.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlPutTask.WriteString("task/")
	urlPutTask.WriteString(taskId)

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)

	resp, err := f.Functions.HttpRequestRetry(http.MethodPut, urlPutTask.String(), f.HttpHeaders, payload, 3)
	if err != nil {
		return errors.New("Error RequestPutTaskStatus request: " + err.Error())
	}

	io.ReadAll(resp.Body)

	return nil
}

func (f *ClickupService) RequestPutTaskStore(taskId string, request type_clickup.TaskRequestStore) error {

	var urlPutTask bytes.Buffer
	urlPutTask.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlPutTask.WriteString("task/")
	urlPutTask.WriteString(taskId)

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)

	resp, err := f.Functions.HttpRequestRetry(http.MethodPut, urlPutTask.String(), f.HttpHeaders, payload, 3)
	if err != nil {
		return errors.New("Error RequestPutTask request: " + err.Error())
	}

	io.ReadAll(resp.Body)

	return nil
}

func (f *ClickupService) RequestSetValueCustomField(taskId string, customFieldId string, request type_clickup.CustomFieldValueRequest) error {

	var urlSetCustomField bytes.Buffer
	urlSetCustomField.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlSetCustomField.WriteString("task/")
	urlSetCustomField.WriteString(taskId)
	urlSetCustomField.WriteString("/field/")
	urlSetCustomField.WriteString(customFieldId)

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)

	resp, err := f.Functions.HttpRequestRetry(http.MethodPost, urlSetCustomField.String(), f.HttpHeaders, payload, 3)
	if err != nil {
		return errors.New("Error RequestSetValueCustomField request: " + err.Error())
	}

	io.ReadAll(resp.Body)

	return nil
}

// func RequestTaskTimeSpent(teamId string, request type_clickup.TaskTimeSpentRequest) error {
// 	var urlTaskTimeSpent bytes.Buffer
// 	urlTaskTimeSpent.WriteString(variables_constant.CLICKUP_API_URL_BASE)
// 	urlTaskTimeSpent.WriteString("team/")
// 	urlTaskTimeSpent.WriteString(teamId)
// 	urlTaskTimeSpent.WriteString("/time_entries")

// 	body, _ := json.Marshal(request)

// 	payload := bytes.NewBuffer(body)

// 	time.Sleep(time.Second)

// 	req, err := http.NewRequest(http.MethodPost, urlTaskTimeSpent.String(), payload)
// 	if err != nil {
// 		return errors.New("Error RequestTaskTimeSpent request: " + err.Error())
// 	}

// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
// 	client := &http.Client{Timeout: time.Second * 10}
// 	resp, err := client.Do(req)
// 	defer req.Body.Close()

// 	if resp.StatusCode != 200 {
// 		return errors.New("Error RequestTaskTimeSpent StatusCode: " + http.StatusText(resp.StatusCode))
// 	}

// 	if err != nil {
// 		return errors.New("Error RequestTaskTimeSpent response: " + err.Error())
// 	}

// 	io.ReadAll(resp.Body)

// 	return nil
// }

func (f *ClickupService) TaskCreateRequest(request type_clickup.TaskCreateRequest) (type_clickup.TaskResponse, error) {
	var result type_clickup.TaskResponse
	var urlCreateTask bytes.Buffer
	urlCreateTask.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlCreateTask.WriteString("list/")
	urlCreateTask.WriteString(variables_global.Customer.ClickUpListId)
	urlCreateTask.WriteString("/task")

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)

	resp, err := f.Functions.HttpRequestRetry(http.MethodPost, urlCreateTask.String(), f.HttpHeaders, payload, 3)
	if err != nil {
		return result, errors.New("Error TaskCreateRequest Request: " + err.Error())
	}

	data, _ := io.ReadAll(resp.Body)

	json.Unmarshal([]byte(string(data)), &result)

	return result, nil
}

func (f *ClickupService) CheckTags(tags []type_clickup.TagResponse) bool {
	ret := false

	for i := 0; i < len(tags); i++ {
		for j := 0; j < len(variables_global.Config.Tags); j++ {
			if strings.EqualFold(tags[i].Name, variables_global.Config.Tags[j].Value) {
				return true
			}
		}
	}

	return ret
}

func (f *ClickupService) CheckSpecificTag(tags []type_clickup.TagResponse, tagVerify string) bool {
	ret := false

	for i := 0; i < len(tags); i++ {
		if strings.EqualFold(tags[i].Name, tagVerify) {
			return true
		}
	}

	return ret
}

func (f *ClickupService) RetDeliveryPointTag(tags []type_clickup.TagResponse) int {
	ret := 0

	for i := 0; i < len(tags); i++ {
		for j := 0; j < len(variables_global.Config.Tags); j++ {
			if strings.EqualFold(tags[i].Name, variables_global.Config.Tags[j].Value) {
				return variables_global.Config.Tags[j].DeliveryPoints
			}
		}
	}

	return ret
}

func (f *ClickupService) VerifyErrorsProjectWithStore(list type_clickup.ListResponse) {
	f.VerifySubtask(list, int(enum_clickup_type_ps_hierarchy.EPIC), int(enum_clickup_type_ps_hierarchy.STORE))
	f.VerifySubtask(list, int(enum_clickup_type_ps_hierarchy.STORE), int(enum_clickup_type_ps_hierarchy.TASK))
	f.VerifyTasks(list)
}

func (f *ClickupService) VerifySubtask(list type_clickup.ListResponse, customFieldTypeConsulting int, customFieldTypeConsultingSubTask int) {

	page := 0

	for {

		tasks, err := f.ReturnTasks(list.Id,
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
			task, err := f.ReturnTask(tasks.Tasks[i].Id)

			if err != nil {
				fmt.Println("Error VerifySubtask GetTask :: ", err.Error())
				return
			}

			if strings.ToLower(task.Status.Status) != "backlog" {

				if strings.EqualFold(task.Parent, "") && customFieldTypeConsulting != int(enum_clickup_type_ps_hierarchy.EPIC) {
					fmt.Println("Store  Without EPIC",
						" :: ", variables_global.Customer.IntegrationName, " :: ", task.Name,
						" :: ", strings.ToLower(task.Status.Status), " :: ", task.Url,
						" :: ", f.RetAssigness(task.Assignees))
					continue
				}

				if len(task.SubTasks) == 0 {
					fmt.Println(enum_clickup_type_ps_hierarchy.ToString(customFieldTypeConsulting),
						" Without ",
						enum_clickup_type_ps_hierarchy.ToString(customFieldTypeConsultingSubTask),
						" :: ", variables_global.Customer.IntegrationName, " :: ", task.Name,
						" :: ", strings.ToLower(task.Status.Status), " :: ", task.Url,
						" :: ", f.RetAssigness(task.Assignees))
					continue
				}

				if variables_global.Customer.ValidatePSTeam && len(task.CustomField.PSTeam) == 0 {
					fmt.Println("EPIC or Story without PS-TEAM: ", variables_global.Customer.IntegrationName, " :: ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", task.Url,
						" :: ", f.RetAssigness(task.Assignees))
				}

				if variables_global.Customer.ValidatePSCustomer && len(task.CustomField.PSCustomer) == 0 {
					fmt.Println("EPIC or Story without PS-Customer: ", variables_global.Customer.IntegrationName, " :: ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", task.Url,
						" :: ", f.RetAssigness(task.Assignees))
				}

				if customFieldTypeConsulting == int(enum_clickup_type_ps_hierarchy.STORE) && variables_global.Customer.ValidateTag {
					if !f.CheckTags(task.Tags) {
						fmt.Println("Story without TAGS", " :: ", variables_global.Customer.IntegrationName, " :: ", task.Name, " :: ",
							strings.ToLower(task.Status.Status), " :: ", task.Url,
							" :: ", f.RetAssigness(task.Assignees))

					}

					if variables_global.Customer.ValidatePSConvisoPlatformLink && (task.CustomField.PSConvisoPlatformLink == "" || !strings.Contains(task.CustomField.PSConvisoPlatformLink, "/projects/")) {
						fmt.Println("Story without Conviso Platform URL: ", " :: ", variables_global.Customer.IntegrationName,
							" :: ", task.Name, " :: ",
							strings.ToLower(task.Status.Status), " :: ", task.Url,
							" :: ", f.RetAssigness(task.Assignees))
					}
				}

				for j := 0; j < len(task.SubTasks); j++ {
					subTask, err := f.ReturnTask(task.SubTasks[j].Id)
					if err != nil {
						fmt.Println("Error VerifySubtask GetTask GetSubTask :: ", err.Error())
						return
					}

					customFieldsSubTask := f.RetCustomFieldTypeConsulting(subTask.CustomFields)

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
							" :: ", subTask.Url, " :: ", f.RetAssigness(subTask.Assignees))
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

func (f *ClickupService) VerifyTasks(list type_clickup.ListResponse) {

	page := 0

	for {

		tasks, err := f.ReturnTasks(list.Id,
			type_clickup.SearchTask{
				TaskType:      int(enum_clickup_type_ps_hierarchy.TASK),
				Page:          page,
				DateUpdatedGt: time.Now().Add(-time.Hour * 240).UTC().UnixMilli(),
				IncludeClosed: true,
				SubTasks:      true,
				TaskStatuses:  "",
			},
		)

		if err != nil {
			fmt.Println("Error VerifyTasks :: ", err.Error())
			return
		}

		for i := 0; i < len(tasks.Tasks); i++ {
			task, err := f.ReturnTask(tasks.Tasks[i].Id)

			if err != nil {
				fmt.Println("Error VerifyTasks GetTask :: ", err.Error())
				return
			}

			if strings.ToLower(task.Status.Status) != "backlog" && !f.CheckSpecificTag(task.Tags, "não executada") {

				if task.Parent == "" {
					fmt.Println("TASK Without Store", " :: ", variables_global.Customer.IntegrationName, " :: ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", task.Url,
						" :: ", f.RetAssigness(task.Assignees))
					continue
				}

				if task.DueDate == "" {
					fmt.Println("Task with errors: ", variables_global.Customer.IntegrationName, " - ", task.Name, " - ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", "DueDate empty", " :: ", task.Url,
						" :: ", f.RetAssigness(task.Assignees))
				}

				if task.StartDate == "" {
					fmt.Println("Task with errors: ", variables_global.Customer.IntegrationName, " - ", task.Name, " - ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", "StartDate empty", " :: ", task.Url,
						" :: ", f.RetAssigness(task.Assignees))
				}

				if task.TimeEstimate == 0 {
					fmt.Println("Task with errors: ", variables_global.Customer.IntegrationName, " - ", task.Name, " - ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", "TimeEstimate empty", " :: ", task.Url,
						" :: ", f.RetAssigness(task.Assignees))
				}

				if (strings.EqualFold(task.Status.Status, "done") || strings.EqualFold(task.Status.Status, "closed")) && task.TimeSpent == 0 {
					fmt.Println("Task with errors: ", variables_global.Customer.IntegrationName, " - ", task.Name, " - ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", "TimeSpent empty", " :: ", task.Url,
						" :: ", f.RetAssigness(task.Assignees))
				}

				if variables_global.Customer.ValidatePSTeam && len(task.CustomField.PSTeam) == 0 {
					fmt.Println("Task with errors: ", variables_global.Customer.IntegrationName, " - ", task.Name, " - ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", "PS-Team empty", " :: ", task.Url,
						" :: ", f.RetAssigness(task.Assignees))
				}

				if variables_global.Customer.ValidatePSCustomer && len(task.CustomField.PSCustomer) == 0 {
					fmt.Println("Task with errors: ", variables_global.Customer.IntegrationName, " - ", task.Name, " - ", task.Name, " :: ",
						strings.ToLower(task.Status.Status), " :: ", "PS-Customer empty", " :: ", task.Url,
						" :: ", f.RetAssigness(task.Assignees))
				}
			}
		}

		if tasks.LastPage {
			break
		}

		page++
	}
}

func (f *ClickupService) UpdateTasksInDoneToClosed(list type_clickup.ListResponse) {
	f.UpdateTasksInDoneToClosedPSHierarchy(list, enum_clickup_type_ps_hierarchy.TASK)
	f.UpdateTasksInDoneToClosedPSHierarchy(list, enum_clickup_type_ps_hierarchy.STORE)
	f.UpdateTasksInDoneToClosedPSHierarchy(list, enum_clickup_type_ps_hierarchy.EPIC)
}

func (f *ClickupService) UpdateTasksInDoneToClosedPSHierarchy(list type_clickup.ListResponse, psHierarchy int) {
	page := 0

	for {
		tasks, err := f.ReturnTasks(list.Id,
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
			err = f.RequestPutTaskStatus(tasks.Tasks[i].Id, type_clickup.TaskRequestStatus{
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

func (f *ClickupService) UpdateProjectWithStore(list type_clickup.ListResponse) {
	f.UpdateTask(list, enum_clickup_type_ps_hierarchy.TASK, enum_clickup_type_ps_hierarchy.STORE)
	f.UpdateTask(list, enum_clickup_type_ps_hierarchy.STORE, enum_clickup_type_ps_hierarchy.EPIC)
}

func (f *ClickupService) UpdateTask(list type_clickup.ListResponse, typeConsultingTask int, typeConsultingParent int) {
	page := 0

	for {

		tasks, err := f.ReturnTasks(list.Id,
			type_clickup.SearchTask{
				TaskType:      typeConsultingTask,
				Page:          page,
				DateUpdatedGt: time.Now().Add(-time.Hour * 240).UTC().UnixMilli(),
				IncludeClosed: true,
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

			taskParent, err := f.ReturnTask((tasks.Tasks[i].Parent))

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
				subTask, err := f.ReturnTask(taskParent.SubTasks[j].Id)
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
				requestTask.Status, hasUpdateStatus = f.RetNewStatus(requestTask.Status, subTask.Status.Status)

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

					projectId, err := f.CPService.RetProjectIdCustomField(taskParent.CustomField.PSConvisoPlatformLink)

					if err == nil {
						convisoPlatformProject, err = f.CPService.GetProject(projectId)
						if err != nil {
							fmt.Println("Error GetProject Conviso Platform :: ", err.Error())
						}
					} else {
						fmt.Println("Error RetProjectIdCustomField Conviso Platform :: ", err.Error())
					}
				}

				if convisoPlatformProject.Id != "" {
					//update the activity in conviso platform project
					err = f.CPService.UpdateActivityRequirement(subTask, convisoPlatformProject)

					if err != nil {
						fmt.Println("Task ", subTask.Name, " not possible update requirement activity in Conviso Platform")
					}
				}
			}

			if allTaskDone {
				requestTask.Status = "closed"
				hasUpdate = true
			}

			if hasUpdate {
				err = f.RequestPutTaskStore(taskParent.Id, requestTask)
				if err != nil {
					fmt.Println("Store not possible update in clickup")
				} else {
					if convisoPlatformProject.Id != "" {
						err = f.CPService.UpdateProjectRest(requestTask, convisoPlatformProject.Id, taskParent.TimeEstimate)
						if err != nil {
							fmt.Println("Store not possible update in conviso platform: " + err.Error())
						}
					}
				}
			}

			if taskParent.CustomField.PSProjectHierarchy == enum_clickup_type_ps_hierarchy.STORE && variables_global.Customer.ValidateTag {
				deliveryPoint := f.RetDeliveryPointTag(taskParent.Tags)
				deliveruPointString := strconv.Itoa(deliveryPoint)
				if deliveryPoint != 0 && !strings.EqualFold(deliveruPointString, taskParent.CustomField.PSDeliveryPoints) {

					err = f.RequestSetValueCustomField(taskParent.Id,
						variables_global.Config.ConfclickUp.CustomFieldPsDeliveryPointsId,
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
