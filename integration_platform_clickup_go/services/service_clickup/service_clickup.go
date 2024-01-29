package service_clickup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"integration_platform_clickup_go/types/type_clickup"
	"integration_platform_clickup_go/utils/functions"
	"integration_platform_clickup_go/utils/variables_constant"
	"integration_platform_clickup_go/utils/variables_global"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var globalClickupHeaders map[string]string

func init() {
	globalClickupHeaders = map[string]string{
		"Authorization": os.Getenv(variables_constant.CLICKUP_TOKEN_NAME),
	}
}

func RetAssigness(assignees []type_clickup.AssigneeField) string {
	ret := ""

	for i := 0; i < len(assignees); i++ {
		ret = ret + assignees[i].Username + ";"
	}

	return ret
}

// func RetCustomerPosition() (string, error) {
// 	result := ""
// 	customFieldsResponse, err := RetCustomFieldCustomerPosition()

// 	if err != nil {
// 		return result, errors.New("Error RetCustomerPosition RequestCustomField: " + err.Error())
// 	}

// 	found := false
// 	for i := 0; i < len(customFieldsResponse.Fields) && found == false; i++ {

// 		if customFieldsResponse.Fields[i].Id == variables_constant.CLICKUP_CUSTOMER_FIELD_ID {
// 			for j := 0; j < len(customFieldsResponse.Fields[i].TypeConfig.Options); j++ {
// 				if strings.ToLower(customFieldsResponse.Fields[i].TypeConfig.Options[j].Name) == strings.ToLower(variables_global.Customer.ClickUpCustomerList) {
// 					result = strconv.Itoa(customFieldsResponse.Fields[i].TypeConfig.Options[j].OrderIndex)
// 					found = true
// 					break
// 				}
// 			}
// 		}
// 	}
// 	return result, nil
// }

func RetTeamPosition(team string) (string, error) {
	result := ""
	customFieldsResponse, err := RetCustomFieldCustomerPosition()

	if err != nil {
		return result, errors.New("Error RetTimePosition RequestCustomField: " + err.Error())
	}

	found := false
	for i := 0; i < len(customFieldsResponse.Fields) && found == false; i++ {

		if customFieldsResponse.Fields[i].Id == variables_constant.CLICKUP_TEAM_FIELD_ID {
			for j := 0; j < len(customFieldsResponse.Fields[i].TypeConfig.Options); j++ {
				if strings.ToLower(customFieldsResponse.Fields[i].TypeConfig.Options[j].Name) == team {
					result = strconv.Itoa(customFieldsResponse.Fields[i].TypeConfig.Options[j].OrderIndex)
					found = true
					break
				}
			}
		}
	}
	return result, nil
}

func RetCustomFieldUrlConviso(customFields []type_clickup.CustomField) string {
	for i := 0; i < len(customFields); i++ {
		if customFields[i].Id == variables_constant.CLICKUP_URL_CONVISO_PLATFORM_FIELD_ID {
			if customFields[i].Value == nil {
				return ""
			} else {
				return customFields[i].Value.(string)
			}
		}
	}
	return ""
}

func RetCustomFieldTypeConsulting(customFields []type_clickup.CustomField) int {
	for i := 0; i < len(customFields); i++ {
		if customFields[i].Id == variables_constant.CLICKUP_TYPE_CONSULTING_FIELD_ID {
			if customFields[i].Value == nil {
				return -1
			} else {
				return int(customFields[i].Value.(float64))
			}
		}
	}
	return 0
}

func RetCustomFieldTeam(customFields []type_clickup.CustomField) []string {
	for i := 0; i < len(customFields); i++ {
		if customFields[i].Id == variables_constant.CLICKUP_TEAM_FIELD_ID {
			if customFields[i].Value == nil {
				return []string{}
			} else {
				aInterface := customFields[i].Value.([]interface{})
				aString := make([]string, len(aInterface))
				for i, v := range aInterface {
					aString[i] = v.(string)
				}
				return aString
			}
		}
	}
	return []string{}
}

func RetCustomFieldCustomerPosition() (type_clickup.CustomFieldsResponse, error) {
	var result type_clickup.CustomFieldsResponse
	var urlGetTasks bytes.Buffer
	urlGetTasks.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlGetTasks.WriteString("list/")
	urlGetTasks.WriteString(variables_global.Customer.ClickUpListId)
	urlGetTasks.WriteString("/field")

	response, err := functions.HttpRequestRetry(http.MethodGet, urlGetTasks.String(), globalClickupHeaders, nil, *variables_global.Config.ConfclickUp.HttpAttempt)

	if err != nil {
		return result, errors.New("Error RetCustomerPosition: " + err.Error())
	}

	data, _ := io.ReadAll(response.Body)

	json.Unmarshal([]byte(string(data)), &result)

	return result, nil
}

func VerifyTasks(taskEpic type_clickup.TaskResponse) error {

	if len(taskEpic.LinkedTasks) == 0 {
		fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", "Nenhuma subtarefa lincada")
	}

	for k := 0; k < len(taskEpic.LinkedTasks); k++ {
		auxTaskId := ""

		if taskEpic.Id == taskEpic.LinkedTasks[k].LinkId {
			auxTaskId = taskEpic.LinkedTasks[k].TaskId
		} else {
			auxTaskId = taskEpic.LinkedTasks[k].LinkId
		}

		taskAux, err := ReturnTask(auxTaskId)
		if err != nil {
			return errors.New("Error taskAux: " + err.Error())
		}

		if strings.ToLower(taskAux.Status.Status) != "backlog" && strings.ToLower(taskAux.Status.Status) != "canceled" && strings.ToLower(taskAux.Status.Status) != "blocked" {
			if taskAux.DueDate == "" {
				fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", taskAux.Name, " :: ", "DueDate empty", " :: ", taskAux.Url)
			}

			if taskAux.StartDate == "" {
				fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", taskAux.Name, " :: ", "StartDate empty", " :: ", taskAux.Url)
			}

			if taskAux.TimeEstimate == 0 {
				fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", taskAux.Name, " :: ", "TimeEstimate empty", " :: ", taskAux.Url)
			}

			if taskAux.Status.Status == "done" && taskAux.TimeSpent == 0 {
				fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", taskAux.Name, " :: ", "TimeSpent empty", " :: ", taskAux.Url)
			}
		}
	}

	return nil
}

func ReturnTasks(listId string, taskType int) (type_clickup.TasksResponse, error) {
	var resultTasks type_clickup.TasksResponse
	var urlGetTasks bytes.Buffer

	intTaskType := strconv.FormatInt(int64(taskType), 10)

	urlGetTasks.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlGetTasks.WriteString("list/")
	urlGetTasks.WriteString(listId)
	urlGetTasks.WriteString("/task?custom_fields=[")
	urlGetTasks.WriteString("{\"field_id\":\"")
	urlGetTasks.WriteString(variables_constant.CLICKUP_TYPE_CONSULTING_FIELD_ID)
	urlGetTasks.WriteString("\",\"operator\":\"=\",\"value\":\"")
	urlGetTasks.WriteString(intTaskType)
	urlGetTasks.WriteString("\"}")
	urlGetTasks.WriteString("]")
	urlGetTasks.WriteString("&include_closed=true")
	urlGetTasks.WriteString("&date_updated_gt=")
	urlGetTasks.WriteString(strconv.FormatInt(time.Now().Add(-time.Hour*240).UTC().UnixMilli(), 10))
	urlGetTasks.WriteString("&subtasks=true")

	response, err := functions.HttpRequestRetry(http.MethodGet, urlGetTasks.String(), globalClickupHeaders, nil, *variables_global.Config.ConfclickUp.HttpAttempt)

	if err != nil {
		return resultTasks, errors.New("Error ReturnTasks: " + err.Error())
	}

	data, _ := io.ReadAll(response.Body)

	json.Unmarshal([]byte(string(data)), &resultTasks)

	return resultTasks, nil
}

func ReturnLists() (type_clickup.ListsResponse, error) {
	var result type_clickup.ListsResponse

	var urlGetLists bytes.Buffer
	urlGetLists.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlGetLists.WriteString("/folder/")
	urlGetLists.WriteString(variables_constant.CLICKUP_FOLDER_CONSULTING_ID)
	urlGetLists.WriteString("/list?archived=false")

	request, err := functions.HttpRequestRetry(http.MethodGet, urlGetLists.String(), globalClickupHeaders, nil, *variables_global.Config.ConfclickUp.HttpAttempt)
	if err != nil {
		return result, errors.New("Error ReturnLists: " + err.Error())
	}

	data, _ := io.ReadAll(request.Body)

	json.Unmarshal([]byte(string(data)), &result)

	return result, nil
}

func ReturnTask(taskId string) (type_clickup.TaskResponse, error) {
	var task type_clickup.TaskResponse
	var urlGetTask bytes.Buffer
	urlGetTask.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlGetTask.WriteString("task/")
	urlGetTask.WriteString(taskId)
	urlGetTask.WriteString("?include_subtasks=true")

	response, err := functions.HttpRequestRetry(http.MethodGet, urlGetTask.String(), globalClickupHeaders, nil, *variables_global.Config.ConfclickUp.HttpAttempt)
	if err != nil {
		return task, errors.New("Error ReturnTask: " + err.Error())
	}

	data, _ := io.ReadAll(response.Body)

	json.Unmarshal([]byte(string(data)), &task)

	//add customFields
	task.CustomField.TypeConsulting = RetCustomFieldTypeConsulting(task.CustomFields)
	task.CustomField.LinkConvisoPlatform = RetCustomFieldUrlConviso(task.CustomFields)
	task.CustomField.Team = RetCustomFieldTeam(task.CustomFields)

	return task, nil
}

func ReturnList(listId string) (type_clickup.ListResponse, error) {
	var list type_clickup.ListResponse
	var urlGetList bytes.Buffer
	urlGetList.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlGetList.WriteString("list/")
	urlGetList.WriteString(listId)
	urlGetList.WriteString("?include_subtasks=true")

	response, err := functions.HttpRequestRetry(http.MethodGet, urlGetList.String(), globalClickupHeaders, nil, *variables_global.Config.ConfclickUp.HttpAttempt)

	if err != nil {
		return list, errors.New("Error ReturnList: " + err.Error())
	}

	data, _ := io.ReadAll(response.Body)

	json.Unmarshal([]byte(string(data)), &list)

	return list, nil
}

func RetNewStatus(statusTask string, statusSubTask string) (string, bool) {

	newReturn := statusTask
	hasUpdate := false

	switch statusSubTask {
	case "backlog":
		break
	case "to do":
		if statusTask == "backlog" {
			newReturn = "to do"
			hasUpdate = true
		}
		break
	case "in progress", "done":
		if statusTask == "backlog" || statusTask == "to do" || statusTask == "blocked" {
			newReturn = "in progress"
			hasUpdate = true
		}
	}
	return newReturn, hasUpdate
}

func RequestPutTask(taskId string, request type_clickup.TaskRequest) error {

	var urlPutTask bytes.Buffer
	urlPutTask.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlPutTask.WriteString("task/")
	urlPutTask.WriteString(taskId)

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)

	time.Sleep(time.Second)

	req, err := http.NewRequest(http.MethodPut, urlPutTask.String(), payload)
	if err != nil {
		return errors.New("Error RequestPutTask request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer req.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Error RequestPutTask StatusCode: " + http.StatusText(resp.StatusCode))
	}

	if err != nil {
		return errors.New("Error RequestPutTask response: " + err.Error())
	}

	io.ReadAll(resp.Body)

	return nil
}

func RequestPutTaskStore(taskId string, request type_clickup.TaskRequestStore) error {

	var urlPutTask bytes.Buffer
	urlPutTask.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlPutTask.WriteString("task/")
	urlPutTask.WriteString(taskId)

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)

	time.Sleep(time.Second)

	req, err := http.NewRequest(http.MethodPut, urlPutTask.String(), payload)
	if err != nil {
		return errors.New("Error RequestPutTask request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer req.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Error RequestPutTask StatusCode: " + http.StatusText(resp.StatusCode))
	}

	if err != nil {
		return errors.New("Error RequestPutTask response: " + err.Error())
	}

	io.ReadAll(resp.Body)

	return nil
}
func RequestTaskTimeSpent(teamId string, request type_clickup.TaskTimeSpentRequest) error {
	var urlTaskTimeSpent bytes.Buffer
	urlTaskTimeSpent.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlTaskTimeSpent.WriteString("team/")
	urlTaskTimeSpent.WriteString(teamId)
	urlTaskTimeSpent.WriteString("/time_entries")

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)

	time.Sleep(time.Second)

	req, err := http.NewRequest(http.MethodPost, urlTaskTimeSpent.String(), payload)
	if err != nil {
		return errors.New("Error RequestTaskTimeSpent request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer req.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Error RequestTaskTimeSpent StatusCode: " + http.StatusText(resp.StatusCode))
	}

	if err != nil {
		return errors.New("Error RequestTaskTimeSpent response: " + err.Error())
	}

	io.ReadAll(resp.Body)

	return nil
}

func TaskCreateRequest(request type_clickup.TaskCreateRequest) (type_clickup.TaskResponse, error) {
	var result type_clickup.TaskResponse
	var urlCreateTask bytes.Buffer
	urlCreateTask.WriteString(variables_constant.CLICKUP_API_URL_BASE)
	urlCreateTask.WriteString("list/")
	urlCreateTask.WriteString(variables_global.Customer.ClickUpListId)
	urlCreateTask.WriteString("/task")

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)

	time.Sleep(time.Second)

	req, err := http.NewRequest(http.MethodPost, urlCreateTask.String(), payload)
	if err != nil {
		return result, errors.New("Error TaskCreateRequest Request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)

	if resp.StatusCode != 200 {
		return result, errors.New("Error TaskCreateRequest Status Code: " + strconv.Itoa(resp.StatusCode))
	}

	defer req.Body.Close()
	if err != nil {
		return result, errors.New("Error TaskCreateRequest ClientDo: " + err.Error())
	}

	data, _ := io.ReadAll(resp.Body)

	json.Unmarshal([]byte(string(data)), &result)

	return result, nil
}

func CheckTags(tags []type_clickup.TagResponse, value string) bool {
	ret := false

	vetValue := strings.Split(value, ";")

	for i := 0; i < len(tags); i++ {
		for j := 0; j < len(vetValue); j++ {
			if strings.ToLower(tags[i].Name) == strings.ToLower(vetValue[j]) {
				return true
			}
		}
	}

	return ret
}
