package ServicesClickup

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

	TypesClickup "integration.platform.clickup/types/types_clickup"
	VariablesConstant "integration.platform.clickup/utils/variables_constant"
)

func ClickUpAutomation(justVerify bool) {
	fmt.Println("...Starting ClickUp Automation...")

	lists, err := ReturnLists()

	if err != nil {
		fmt.Println("Error return list: ", err.Error())
		return
	}

	fmt.Println("...Searching valid list...")
	for i := 0; i < len(lists.Lists); i++ {
		if strings.Contains(strings.ToLower(lists.Lists[i].Name), "testeprojetos") {

			fmt.Println("Found valid list ", lists.Lists[i].Name)

			tasks, err := ReturnTasks(lists.Lists[i].Id)

			if err != nil {
				fmt.Println("Error return tasks: ", err.Error())
				return
			}

			for j := 0; j < len(tasks.Tasks); j++ {
				fmt.Println("Task ", j+1, "/", len(tasks.Tasks), " - ", tasks.Tasks[j].Name)

				taskEpic, err := ReturnTask(tasks.Tasks[j].LinkedTasks[0].TaskId)
				if err != nil {
					fmt.Println("Error return task: ", err.Error())
					return
				}

				if justVerify {
					VerifyTasks(taskEpic)
				} else {
					err = ChangeEpicTask(taskEpic, tasks.Tasks[j])

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

func ChangeEpicTask(taskEpic TypesClickup.TaskResponse, taskTask TypesClickup.TaskResponse) error {
	allSubTasksDone := true
	var timeSpent int64
	var requestTask TypesClickup.TaskRequest

	for k := 0; k < len(taskEpic.LinkedTasks); k++ {
		taskAux, err := ReturnTask(taskEpic.LinkedTasks[k].LinkId)
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

		if taskAux.Status.Status != "done" {
			allSubTasksDone = false
		}
	}

	if allSubTasksDone {
		var taskTimeSpentRequest TypesClickup.TaskTimeSpentRequest
		taskTimeSpentRequest.Duration = timeSpent - taskEpic.TimeSpent
		taskTimeSpentRequest.Start = time.Now().UTC().UnixMilli()
		taskTimeSpentRequest.TaskId = taskEpic.Id
		requestTask.Status = "done"
		RequestTaskTimeSpent(taskEpic.TeamId, taskTimeSpentRequest)
	} else {
		requestTask.Status = RetNewStatus(taskEpic.Status.Status, taskTask.Status.Status)
	}

	err := RequestPutTask(taskEpic.Id, requestTask)

	if err != nil {
		return errors.New("Error taskAux: " + err.Error())
	}

	return nil
}

func VerifyTasks(taskEpic TypesClickup.TaskResponse) error {

	for k := 0; k < len(taskEpic.LinkedTasks); k++ {
		taskAux, err := ReturnTask(taskEpic.LinkedTasks[k].LinkId)
		if err != nil {
			return errors.New("Error taskAux: " + err.Error())
		}

		if taskAux.DueDate == "" || taskAux.StartDate == "" || taskAux.TimeEstimate == 0 || (taskAux.Status.Status == "done") && taskAux.TimeSpent == 0 {
			fmt.Println("Task with errors: ", taskEpic.List.Name, " - ", taskEpic.Name, " - ", taskAux.Name)
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
	urlGetTasks.WriteString(strconv.FormatInt(time.Now().Add(-time.Hour*24).UTC().UnixMilli(), 10))

	req, err := http.NewRequest(http.MethodGet, urlGetTasks.String(), nil)
	if err != nil {
		return resultTasks, errors.New("Error request Tasks: " + err.Error())
	}
	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return resultTasks, errors.New("Error response Tasks: " + err.Error())
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
		return result, errors.New("Error request Lists: " + err.Error())
	}

	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return result, errors.New("Error response Lists: " + err.Error())
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
		return task, errors.New("Error request task: " + err.Error())
	}

	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return task, errors.New("Error response task: " + err.Error())
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
		return errors.New("Error request put task: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		return errors.New("Error response put task: " + err.Error())
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
		return errors.New("Error request time spent: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("CLICKUP_TOKEN"))
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer req.Body.Close()
	if err != nil {
		return errors.New("Error response put task: " + err.Error())
	}

	ioutil.ReadAll(resp.Body)

	return nil
}
