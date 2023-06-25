package prompt

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/pkg/model"
	"github.com/pjgaetan/airflow-cli/pkg/utils"
)

func PromptDag() (model.Dag, error) {
	response := request.AirflowGetRequest("dags")
	var dags model.Dags
	if err := json.Unmarshal([]byte(response), &dags); err != nil {
		panic(err)
	}
	dagIds := make([]string, len(dags.Dags))
	for index, dag := range dags.Dags {
		dagIds[index] = dag.DagId
	}
	if len(dags.Dags) == 0 {
		utils.Failed("No run for this dag.")
	}
	var dagId string
	prompt := &survey.Select{
		Message: "Choose a dagRun:",
		Options: dagIds,
	}

	err := survey.AskOne(prompt, &dagId)
	utils.ExitIfError(err)
	if dagId == "" {
		return model.Dag{}, nil
	}
	for _, dag := range dags.Dags {
		if dag.DagId == dagId {
			return dag, nil
		}
	}
	return model.Dag{}, errors.New("No such a dagId error")
}

func PromptDagRun(dagId, orderBy string, limit int) (model.DagRun, error) {
	response := request.AirflowGetRequest("dags/"+dagId+"/dagRuns", [2]string{"limit", strconv.Itoa(limit)}, [2]string{"order_by", orderBy})
	var dagRun model.DagRuns
	if err := json.Unmarshal([]byte(response), &dagRun); err != nil {
		panic(err)
	}
	runs := dagRun.DagRun

	// reverse sort to have last run in first
	sort.Slice(runs, func(i, j int) bool {
		return runs[i].StartDate.Format(time.RFC3339) > runs[j].StartDate.Format(time.RFC3339)
	})

	dagRunIds := make([]string, len(runs))
	for index, run := range runs {
		dagRunIds[index] = run.DagRunId
	}
	if len(runs) == 0 {
		utils.Failed("No run for this dag.")
	}
	fmt.Println(runs[0].DagId)

	var dagRunId string
	prompt := &survey.Select{
		Message: "Choose a dagRun:",
		Options: dagRunIds,
		Description: func(value string, index int) string {
			for _, run := range runs {
				if run.DagRunId == value {
					if run.State == "" {
						return "none"
					}
					return run.State
				}
			}
			return ""
		},
	}

	err := survey.AskOne(prompt, &dagRunId)
	utils.ExitIfError(err)

	if dagRunId == "" {
		return model.DagRun{}, nil
	}

	for _, run := range runs {
		if run.DagRunId == dagRunId {
			return run, nil
		}
	}
	return model.DagRun{}, errors.New("No such a dagRunId error")
}

func PromptSetTaskInstanceState(dagId, dagRunId string) {
	responseTask := request.AirflowGetRequest("dags/" + dagId + "/dagRuns/" + dagRunId + "/taskInstances")
	var taskInstance model.TaskInstances
	if err := json.Unmarshal([]byte(responseTask), &taskInstance); err != nil {
		panic(err)
	}

	instances := taskInstance.TaskInstance
	// reverse sort to have last run in first
	sort.Slice(instances, func(i, j int) bool {
		return instances[i].StartDate.Format(time.RFC3339) > instances[j].StartDate.Format(time.RFC3339)
	})
	taskInstanceIds := make([]string, len(instances))
	for index, task := range instances {
		taskInstanceIds[index] = task.TaskId
	}
	if len(instances) == 0 {
		fmt.Println("No task in this dag for this dag run.")
	}
	fmt.Println(dagRunId)

	// Task selection
	var taskId string
	promptTask := &survey.Select{
		Message: "What task do you want to set?",
		Options: taskInstanceIds,
		Description: func(value string, index int) string {
			for _, i := range instances {
				if i.TaskId == value {
					if i.State == "" {
						return "none"
					}
					return i.State
				}
			}
			return ""
		},
	}
	err := survey.AskOne(promptTask, &taskId)
	utils.ExitIfError(err)
	if taskId == "" {
		return
	}

	// State selection
	var state string
	PromtpState := &survey.Select{
		Message: "What state do you want to state?",
		Options: []string{"success", "failed"},
	}
	err = survey.AskOne(PromtpState, &state)
	utils.ExitIfError(err)

	// Params selections
	params := []string{}

	promptParams := &survey.MultiSelect{
		Message: "Who does it needs to apply to?",
		Options: []string{"include_upstream", "include_downstream"},
	}
	err = survey.AskOne(promptParams, &params)
	utils.ExitIfError(err)
	mapParams := make(map[string]any)
	// Default params
	mapParams["include_future"] = false
	mapParams["include_past"] = false
	mapParams["include_downstream"] = false
	mapParams["include_upstream"] = false
	mapParams["dag_run_id"] = dagRunId

	for _, p := range params {
		mapParams[p] = true
	}
	mapParams["new_state"] = state
	mapParams["task_id"] = taskId
	mapParams["dry_run"] = false
	jsonParams, err := json.Marshal(mapParams)
	if err != nil {
		log.Fatal("Error ", err)
	}

	response := request.AirflowPostRequest("dags/"+dagId+"/updateTaskInstancesState", string(jsonParams))
	var dagRun model.DagRuns
	if err := json.Unmarshal([]byte(response), &dagRun); err != nil {
		panic(err)
	}
	fmt.Println(string(response))
}
