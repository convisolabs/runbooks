package ServiceClickup

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
	"time"

	"golang.org/x/exp/slices"
	TypesClickup "integration.platform.clickup/types/type_clickup"
	Functions "integration.platform.clickup/utils/functions"
	VariablesConstant "integration.platform.clickup/utils/variables_constant"
	VariablesGlobal "integration.platform.clickup/utils/variables_global"
)

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

func RetCustomFieldCustomerPosition() (TypesClickup.CustomFieldsResponse, error) {
	var result TypesClickup.CustomFieldsResponse
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

func ClickUpAutomation(justVerify bool) {
	fmt.Println("...Starting ClickUp Automation...")

	lists, err := ReturnLists()

	if err != nil {
		fmt.Println("Error return list: ", err.Error())
		return
	}

	fmt.Println("...Searching valid list...")
	for i := 0; i < len(lists.Lists); i++ {

		fmt.Println("Found List ", lists.Lists[i].Name)

		if Functions.CustomerExistsYamlFileByClickUpListId(lists.Lists[i].Id, Functions.LoadCustomerByYamlFile()) {

			var sliceEpicId []string

			fmt.Println("Found valid list ", lists.Lists[i].Name)

			time.Sleep(time.Second)

			tasks, err := ReturnTasks(lists.Lists[i].Id)

			if err != nil {
				fmt.Println("Error return tasks: ", err.Error())
				return
			}

			for j := 0; j < len(tasks.Tasks); j++ {
				fmt.Println("Task ", j+1, "/", len(tasks.Tasks), " - ", tasks.Tasks[j].Name)

				auxEpicTaskId := ""

				if len(tasks.Tasks[j].LinkedTasks) == 0 {
					fmt.Println("Error 0 epics", " :: ", lists.Lists[i].Name, " - ", tasks.Tasks[j].Name)
					continue
				}

				if len(tasks.Tasks[j].LinkedTasks) > 1 {
					fmt.Println("Error 2 epics:", " :: ", lists.Lists[i].Name, " - ", tasks.Tasks[j].Name)
					continue
				}

				//dependendo a ordem que vc linkar as tarefas ele vai jogar no linkid ou no taskid
				if tasks.Tasks[j].Id == tasks.Tasks[j].LinkedTasks[0].TaskId {
					auxEpicTaskId = tasks.Tasks[j].LinkedTasks[0].LinkId
				} else {
					auxEpicTaskId = tasks.Tasks[j].LinkedTasks[0].TaskId
				}

				if slices.Contains(sliceEpicId, auxEpicTaskId) {
					continue
				}

				sliceEpicId = append(sliceEpicId, auxEpicTaskId)

				time.Sleep(time.Second)

				taskEpic, err := ReturnTask(auxEpicTaskId)
				if err != nil {
					fmt.Println("Error return task: ", err.Error())
					return
				}

				if justVerify {
					time.Sleep(time.Second)
					VerifyTasks(taskEpic)
				} else {
					time.Sleep(time.Second)
					err = ChangeEpicTask(taskEpic)

					if err != nil {
						fmt.Println("Error change Epic Task: ", err.Error())
						return
					}
				}
			}
		}
	}
	fmt.Println("...Finishing ClickUp Automation...")
}

func ChangeEpicTask(taskEpic TypesClickup.TaskResponse) error {
	allSubTasksDone := true
	var timeSpent int64
	var requestTask TypesClickup.TaskRequest

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
		auxDuoDate, _ := strconv.ParseInt(taskAux.DueDate, 10, 64)
		auxStartDate, _ := strconv.ParseInt(taskAux.StartDate, 10, 64)
		requestTask.TimeEstimate = requestTask.TimeEstimate + taskAux.TimeEstimate
		timeSpent = timeSpent + taskAux.TimeSpent
		if auxDuoDate > requestTask.DueDate {
			requestTask.DueDate = auxDuoDate
		}

		if auxStartDate != 0 && auxStartDate < requestTask.StartDate || requestTask.StartDate == 0 {
			requestTask.StartDate = auxStartDate
		}

		if taskAux.Status.Status != "done" && taskAux.Status.Status != "canceled" {
			allSubTasksDone = false
		}

		requestTask.Status = RetNewStatus(taskEpic.Status.Status, taskAux.Status.Status)
		taskEpic.Status.Status = requestTask.Status
	}

	if allSubTasksDone {
		requestTask.Status = "done"
	}

	if (timeSpent - taskEpic.TimeSpent) > 0 {
		var taskTimeSpentRequest TypesClickup.TaskTimeSpentRequest
		taskTimeSpentRequest.Duration = timeSpent - taskEpic.TimeSpent
		taskTimeSpentRequest.Start = time.Now().UTC().UnixMilli()
		taskTimeSpentRequest.TaskId = taskEpic.Id
		RequestTaskTimeSpent(taskEpic.TeamId, taskTimeSpentRequest)
	}

	err := RequestPutTask(taskEpic.Id, requestTask)

	if err != nil {
		return errors.New("Error taskAux: " + err.Error())
	}

	return nil
}

func VerifyTasks(taskEpic TypesClickup.TaskResponse) error {

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
				fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", taskAux.Name, " :: ", "DueDate empty")
			}

			if taskAux.StartDate == "" {
				fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", taskAux.Name, " :: ", "StartDate empty")
			}

			if taskAux.TimeEstimate == 0 {
				fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", taskAux.Name, " :: ", "TimeEstimate empty")
			}

			if taskAux.Status.Status == "done" && taskAux.TimeSpent == 0 {
				fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", taskAux.Name, " :: ", "TimeSpent empty")
			}
		}
	}

	return nil
}

func ReturnTasks(listId string) (TypesClickup.TasksResponse, error) {
	var resultTasks TypesClickup.TasksResponse
	var urlGetTasks bytes.Buffer
	urlGetTasks.WriteString(VariablesConstant.CLICKUP_API_URL_BASE)
	urlGetTasks.WriteString("list/")
	urlGetTasks.WriteString(listId)
	urlGetTasks.WriteString("/task?custom_fields=[")
	urlGetTasks.WriteString("{\"field_id\":\"664816bc-a899-45ec-9801-5a1e5be9c5f6\",\"operator\":\">=\",\"value\":\"1\"}")
	urlGetTasks.WriteString("]")
	urlGetTasks.WriteString("&include_closed=true")
	urlGetTasks.WriteString("&date_updated_gt=")
	urlGetTasks.WriteString(strconv.FormatInt(time.Now().Add(-time.Hour*96).UTC().UnixMilli(), 10))

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

	data, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(string(data)), &resultTasks)
	return resultTasks, nil
}

func ReturnLists() (TypesClickup.ListsResponse, error) {
	var result TypesClickup.ListsResponse

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

	data, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal([]byte(string(data)), &result)

	return result, nil
}

func ReturnTask(taskId string) (TypesClickup.TaskResponse, error) {
	var task TypesClickup.TaskResponse
	var urlGetTask bytes.Buffer
	urlGetTask.WriteString(VariablesConstant.CLICKUP_API_URL_BASE)
	urlGetTask.WriteString("task/")
	urlGetTask.WriteString(taskId)

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
	data, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal([]byte(string(data)), &task)

	return task, nil
}

func RetNewStatus(statusPrincipal string, statusLinked string) string {

	newReturn := statusPrincipal

	switch statusLinked {
	case "backlog":
		break
	case "to do":
		if statusPrincipal == "backlog" {
			newReturn = "to do"
		}
		break
	case "in progress", "done":
		if statusPrincipal == "backlog" || statusPrincipal == "to do" || statusPrincipal == "blocked" {
			newReturn = "in progress"
		}
	}
	return newReturn
}

func RequestPutTask(taskId string, request TypesClickup.TaskRequest) error {

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

func RequestTaskTimeSpent(teamId string, request TypesClickup.TaskTimeSpentRequest) error {
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

func TaskCreateRequest(request TypesClickup.TaskCreateRequest) (TypesClickup.TaskResponse, error) {
	var result TypesClickup.TaskResponse
	var urlCreateTask bytes.Buffer
	urlCreateTask.WriteString(VariablesConstant.CLICKUP_API_URL_BASE)
	urlCreateTask.WriteString("list/")
	urlCreateTask.WriteString(VariablesGlobal.Customer.ClickUpListId)
	urlCreateTask.WriteString("/task")

	fmt.Println(urlCreateTask.String())

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
