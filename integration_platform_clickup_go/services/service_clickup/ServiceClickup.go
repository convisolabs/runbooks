package ServiceClickup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	TypeClickup "integration.platform.clickup/types/type_clickup"
	VariablesConstant "integration.platform.clickup/utils/variables_constant"
	VariablesGlobal "integration.platform.clickup/utils/variables_global"
)

func RetAssigness(assignees []TypeClickup.AssigneeField) string {
	ret := ""

	for i := 0; i < len(assignees); i++ {
		ret = ret + assignees[i].Username + ";"
	}

	return ret
}

func RetCustomerPosition() (string, error) {
	result := ""
	customFieldsResponse, err := RetCustomFieldCustomerPosition()

	if err != nil {
		return result, errors.New("Error RetCustomerPosition RequestCustomField: " + err.Error())
	}

	found := false
	for i := 0; i < len(customFieldsResponse.Fields) && found == false; i++ {

		if customFieldsResponse.Fields[i].Id == VariablesConstant.CLICKUP_CUSTOMER_FIELD_ID {
			for j := 0; j < len(customFieldsResponse.Fields[i].TypeConfig.Options); j++ {
				if strings.ToLower(customFieldsResponse.Fields[i].TypeConfig.Options[j].Name) == strings.ToLower(VariablesGlobal.Customer.ClickUpCustomerList) {
					result = strconv.Itoa(customFieldsResponse.Fields[i].TypeConfig.Options[j].OrderIndex)
					found = true
					break
				}
			}
		}
	}
	return result, nil
}

func RetCustomFieldUrlConviso(customFields []TypeClickup.CustomField) string {
	for i := 0; i < len(customFields); i++ {
		if customFields[i].Id == VariablesConstant.CLICKUP_URL_CONVISO_PLATFORM_FIELD_ID {
			return customFields[i].Value.(string)
		}
	}
	return ""
}

func RetCustomFieldTypeConsulting(customFields []TypeClickup.CustomField) int {
	for i := 0; i < len(customFields); i++ {
		if customFields[i].Id == VariablesConstant.CLICKUP_TYPE_CONSULTING_FIELD_ID {
			return int(customFields[i].Value.(float64))
		}
	}
	return 0
}

func RetCustomFieldCustomerPosition() (TypeClickup.CustomFieldsResponse, error) {
	var result TypeClickup.CustomFieldsResponse
	var urlGetTasks bytes.Buffer
	urlGetTasks.WriteString(VariablesConstant.CLICKUP_API_URL_BASE)
	urlGetTasks.WriteString("list/")
	urlGetTasks.WriteString(VariablesGlobal.Customer.ClickUpListId)
	urlGetTasks.WriteString("/field")

	req, err := http.NewRequest(http.MethodGet, urlGetTasks.String(), nil)
	if err != nil {
		return result, errors.New("Error RetCustomerPosition NewRequest: " + err.Error())
	}
	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		return result, errors.New("Error RetCustomerPosition StatusCode: " + http.StatusText(resp.StatusCode))
	}

	if err != nil {
		return result, errors.New("Error RetCustomerPosition ClientDo: " + err.Error())
	}
	data, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(string(data)), &result)

	return result, nil
}

func VerifyTasks(taskEpic TypeClickup.TaskResponse) error {

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

func ReturnTasks(listId string, taskType int) (TypeClickup.TasksResponse, error) {
	var resultTasks TypeClickup.TasksResponse
	var urlGetTasks bytes.Buffer

	teste := strconv.FormatInt(int64(taskType), 10)

	urlGetTasks.WriteString(VariablesConstant.CLICKUP_API_URL_BASE)
	urlGetTasks.WriteString("list/")
	urlGetTasks.WriteString(listId)
	urlGetTasks.WriteString("/task?custom_fields=[")
	urlGetTasks.WriteString("{\"field_id\":\"")
	urlGetTasks.WriteString(VariablesConstant.CLICKUP_TYPE_CONSULTING_FIELD_ID)
	urlGetTasks.WriteString("\",\"operator\":\"=\",\"value\":\"")
	urlGetTasks.WriteString(teste)
	urlGetTasks.WriteString("\"}")
	urlGetTasks.WriteString("]")
	urlGetTasks.WriteString("&include_closed=true")
	urlGetTasks.WriteString("&date_updated_gt=")
	urlGetTasks.WriteString(strconv.FormatInt(time.Now().Add(-time.Hour*240).UTC().UnixMilli(), 10))

	if VariablesGlobal.Customer.HasStore {
		urlGetTasks.WriteString("&subtasks=true")
	}

	req, err := http.NewRequest(http.MethodGet, urlGetTasks.String(), nil)
	if err != nil {
		return resultTasks, errors.New("Error ReturnTasks request: " + err.Error())
	}
	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)

	if resp.StatusCode != 200 {
		return resultTasks, errors.New("Error ReturnTasks StatusCode: " + http.StatusText(resp.StatusCode))
	}

	if err != nil {
		return resultTasks, errors.New("Error ReturnTasks response: " + err.Error())
	}

	data, _ := io.ReadAll(resp.Body)
	json.Unmarshal([]byte(string(data)), &resultTasks)
	return resultTasks, nil
}

func ReturnLists() (TypeClickup.ListsResponse, error) {
	var result TypeClickup.ListsResponse

	var urlGetLists bytes.Buffer
	urlGetLists.WriteString(VariablesConstant.CLICKUP_API_URL_BASE)
	urlGetLists.WriteString("/folder/")
	urlGetLists.WriteString(VariablesConstant.CLICKUP_FOLDER_CONSULTING_ID)
	urlGetLists.WriteString("/list?archived=false")

	req, err := http.NewRequest(http.MethodGet, urlGetLists.String(), nil)
	if err != nil {
		return result, errors.New("Error ReturnLists request: " + err.Error())
	}

	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)

	if resp.StatusCode != 200 {
		return result, errors.New("Error ReturnLists StatusCode: " + http.StatusText(resp.StatusCode))
	}

	if err != nil {
		return result, errors.New("Error ReturnLists response: " + err.Error())
	}

	data, _ := io.ReadAll(resp.Body)

	json.Unmarshal([]byte(string(data)), &result)

	return result, nil
}

func ReturnTask(taskId string) (TypeClickup.TaskResponse, error) {
	var task TypeClickup.TaskResponse
	var urlGetTask bytes.Buffer
	urlGetTask.WriteString(VariablesConstant.CLICKUP_API_URL_BASE)
	urlGetTask.WriteString("task/")
	urlGetTask.WriteString(taskId)

	if VariablesGlobal.Customer.HasStore {
		urlGetTask.WriteString("?include_subtasks=true")
	}

	req, err := http.NewRequest(http.MethodGet, urlGetTask.String(), nil)
	if err != nil {
		return task, errors.New("Error ReturnTask request: " + err.Error())
	}

	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)

	if resp.StatusCode != 200 {
		return task, errors.New("Error ReturnTask StatusCode: " + http.StatusText(resp.StatusCode))
	}

	if err != nil {
		return task, errors.New("Error ReturnTask response: " + err.Error())
	}
	data, _ := io.ReadAll(resp.Body)

	json.Unmarshal([]byte(string(data)), &task)

	return task, nil
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

func RequestPutTask(taskId string, request TypeClickup.TaskRequest) error {

	var urlPutTask bytes.Buffer
	urlPutTask.WriteString(VariablesConstant.CLICKUP_API_URL_BASE)
	urlPutTask.WriteString("task/")
	urlPutTask.WriteString(taskId)

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)
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

	ioutil.ReadAll(resp.Body)

	return nil
}

func RequestPutTaskStore(taskId string, request TypeClickup.TaskRequestStore) error {

	var urlPutTask bytes.Buffer
	urlPutTask.WriteString(VariablesConstant.CLICKUP_API_URL_BASE)
	urlPutTask.WriteString("task/")
	urlPutTask.WriteString(taskId)

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)
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
func RequestTaskTimeSpent(teamId string, request TypeClickup.TaskTimeSpentRequest) error {
	var urlTaskTimeSpent bytes.Buffer
	urlTaskTimeSpent.WriteString(VariablesConstant.CLICKUP_API_URL_BASE)
	urlTaskTimeSpent.WriteString("team/")
	urlTaskTimeSpent.WriteString(teamId)
	urlTaskTimeSpent.WriteString("/time_entries")

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)
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

	ioutil.ReadAll(resp.Body)

	return nil
}

func TaskCreateRequest(request TypeClickup.TaskCreateRequest) (TypeClickup.TaskResponse, error) {
	var result TypeClickup.TaskResponse
	var urlCreateTask bytes.Buffer
	urlCreateTask.WriteString(VariablesConstant.CLICKUP_API_URL_BASE)
	urlCreateTask.WriteString("list/")
	urlCreateTask.WriteString(VariablesGlobal.Customer.ClickUpListId)
	urlCreateTask.WriteString("/task")

	body, _ := json.Marshal(request)

	payload := bytes.NewBuffer(body)
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

	data, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal([]byte(string(data)), &result)

	return result, nil
}
