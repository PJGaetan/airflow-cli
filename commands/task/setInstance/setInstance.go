package setInstance

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/pkg/model"
	"github.com/pjgaetan/airflow-cli/pkg/prompt"

	"github.com/spf13/cobra"
)

var (
	DagId    string
	DagRunId string
	Limit    int
	OrderBy  string
)

// listCmd represents the list command
func NewSet() *cobra.Command {
	stateCmd := cobra.Command{
		Use:   "set-instance",
		Short: "Set task instances state (success,failed)",
		Run:   state,
	}
	stateCmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	stateCmd.Flags().StringVarP(&DagRunId, "dag-run-id", "r", "", "dag id")
	stateCmd.Flags().StringVarP(&OrderBy, "order-by", "o", "-start_date", "The name of the field to order the results by. Prefix a field name with - to reverse the sort order.")
	stateCmd.Flags().IntVarP(&Limit, "limit", "l", 20, "The numbers of items to return.")
	return &stateCmd
}

func state(cmd *cobra.Command, args []string) {
	if DagId == "" {
		dag, err := prompt.PromptDag()
		if err != nil {
			panic(err)
		}
		if reflect.DeepEqual(dag, model.Dag{}) {
			os.Exit(0)
		}
		DagId = dag.Dag_id
	}

	if DagRunId == "" {
		run, err := prompt.PromptDagRun(DagId, OrderBy, Limit)
		if err != nil {
			panic(err)
		}
		if reflect.DeepEqual(run, model.Dag{}) {
			os.Exit(0)
		}
		DagRunId = run.Dag_run_id
	}

	responseTask := request.AirflowGetRequest("dags/" + DagId + "/dagRuns/" + DagRunId + "/taskInstances")
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
	fmt.Println(DagRunId)

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
	survey.AskOne(promptTask, &taskId)
	if taskId == "" {
		return
	}

	// State selection
	var state string
	PromtpState := &survey.Select{
		Message: "What state do you want to state?",
		Options: []string{"success", "failed"},
	}
	survey.AskOne(PromtpState, &state)

	// Params selections
	params := []string{}

	promptParams := &survey.MultiSelect{
		Message: "Who does it needs to apply to?",
		Options: []string{"include_upstream", "include_downstream"},
	}
	survey.AskOne(promptParams, &params)
	mapParams := make(map[string]any)
	// Default params
	mapParams["include_future"] = false
	mapParams["include_past"] = false
	mapParams["include_downstream"] = false
	mapParams["include_upstream"] = false
	mapParams["dag_run_id"] = DagRunId

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

	response := request.AirflowPostRequest("dags/"+DagId+"/updateTaskInstancesState", string(jsonParams))
	var dagRun model.DagRuns
	if err := json.Unmarshal([]byte(response), &dagRun); err != nil {
		panic(err)
	}
	fmt.Println(string(response))
}
