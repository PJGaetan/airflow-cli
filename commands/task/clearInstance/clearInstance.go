package clearInstance

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
	"github.com/pjgaetan/airflow-cli/internal/constant"
	"github.com/pjgaetan/airflow-cli/pkg/model"
	"github.com/pjgaetan/airflow-cli/pkg/prompt"
	"github.com/pjgaetan/airflow-cli/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	DagId    string
	DagRunId string
	Limit    int
	OrderBy  string
)

// NewClear represents the list command.
func NewClear() *cobra.Command {
	clearCmd := cobra.Command{
		Use:   "clear-instance",
		Short: "Clear task instances",
		Run:   clear,
	}
	clearCmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	clearCmd.Flags().StringVarP(&DagRunId, "dag-run-id", "r", "", "dag id")
	clearCmd.Flags().StringVarP(&OrderBy, "order-by", "o", "-start_date", "The name of the field to order the results by. Prefix a field name with - to reverse the sort order.")
	clearCmd.Flags().IntVarP(&Limit, "limit", "l", constant.DEAULT_ITEM_LIMIT, "The numbers of items to return.")
	return &clearCmd
}

func clear(cmd *cobra.Command, args []string) {
	if DagId == "" {
		dag, err := prompt.PromptDag()
		if err != nil {
			panic(err)
		}
		if reflect.DeepEqual(dag, model.Dag{}) {
			os.Exit(0)
		}
		DagId = dag.DagId
	}

	if DagRunId == "" {
		run, err := prompt.PromptDagRun(DagId, OrderBy, Limit)
		if err != nil {
			panic(err)
		}
		if reflect.DeepEqual(run, model.Dag{}) {
			os.Exit(0)
		}
		DagRunId = run.DagRunId
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
	var taskId string
	promptTask := &survey.Select{
		Message: "What task do you want to clear?",
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
	var task model.TaskInstance
	for _, t := range instances {
		if t.TaskId == taskId {
			task = t
			break

		}
	}
	if reflect.DeepEqual(task, model.TaskInstance{}) {
		utils.Failed("no such task as " + taskId)
	}

	params := []string{}
	promptParams := &survey.MultiSelect{
		Message: "What parameters to apply?",
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
	mapParams["only_failed"] = false
	mapParams["dry_run"] = false

	for _, p := range params {
		mapParams[p] = true
	}

	mapParams["task_ids"] = []string{taskId}

	// Must provide either dagrunid or start/end dates
	mapParams["dag_run_id"] = DagRunId
	// mapParams["start_date"] = task.StartDate.Format(time.RFC3339)
	// mapParams["end_date"] = task.EndDate.Format(time.RFC3339)

	jsonParams, err := json.Marshal(mapParams)
	if err != nil {
		log.Fatal("Error ", err)
	}

	response := request.AirflowPostRequest("dags/"+DagId+"/clearTaskInstances", string(jsonParams))
	var dagRun model.DagRuns
	if err := json.Unmarshal([]byte(response), &dagRun); err != nil {
		panic(err)
	}
	fmt.Println(string(response))
}
