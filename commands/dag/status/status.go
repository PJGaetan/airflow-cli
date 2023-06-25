package status

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pjgaetan/airflow-cli/api/request"
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

const DEFAULT_RETURN_ITEMS = 5

// NewStatus represents the list command.
func NewStatus() *cobra.Command {
	statusCmd := cobra.Command{
		Use:   "status",
		Short: "Status of dag runs (tasks & logs)",
		Run:   status,
	}
	statusCmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	statusCmd.Flags().IntVarP(&Limit, "limit", "l", DEFAULT_RETURN_ITEMS, "The numbers of items to return.")
	statusCmd.Flags().StringVarP(&OrderBy, "order-by", "o", "-start_date", "The name of the field to order the results by. Prefix a field name with - to reverse the sort order.")
	return &statusCmd
}

func showLog() {
	responseTask := request.AirflowGetRequest("dags/" + DagId + "/dagRuns/" + DagRunId + "/taskInstances")
	var taskInstance model.TaskInstances
	if err := json.Unmarshal([]byte(responseTask), &taskInstance); err != nil {
		panic(err)
	}

	instances := taskInstance.TaskInstance
	// reverse sort to have last run in first
	sort.Slice(instances, func(i, j int) bool {
		return instances[i].QueuedWhen.Format(time.RFC3339) < instances[j].QueuedWhen.Format(time.RFC3339)
	})

	taskInstanceIds := make([]string, len(instances))
	for index, instance := range instances {
		taskInstanceIds[index] = instance.TaskId
	}

	var taskInstanceId string
	promptLogs := &survey.Select{
		Message: "Choose a task instance:",
		Options: taskInstanceIds,
		Description: func(value string, index int) string {
			for _, instance := range instances {
				if instance.TaskId == value {
					var state string
					if instance.State == "" {
						state = "none"
					} else {
						state = instance.State
					}
					return state + " (" + strconv.Itoa(int(instance.TryNumber)) + ") - " + instance.QueuedWhen.Format(time.RFC3339)
				}
			}
			return ""
		},
	}
	err := survey.AskOne(promptLogs, &taskInstanceId)
	utils.ExitIfError(err)
	responseLogs := request.AirflowGetRequest("dags/"+DagId+"/dagRuns/"+DagRunId+"/taskInstances/"+taskInstanceId+"/logs/"+"1", [2]string{"full_content", "true"})
	var logs model.Logs
	if err := json.Unmarshal([]byte(responseLogs), &logs); err != nil {
		panic(err)
	}

	log := strings.Trim(string(logs.Content), "[]()\"")
	for _, line := range strings.Split(log, "\\n") {
		line = strings.ReplaceAll(line, "\\\\", "\\")
		line = strings.Trim(line, "\\")
		fmt.Println(line)
	}
}

func setDagRunState(dagRunState string) {
	state := ""
	prompt := &survey.Select{
		Message: "Dag " + DagId + " of run '" + DagRunId + "' is in state: " + dagRunState + "\nSet dagRun to which state ?",
		Options: []string{"success", "failed", "queued"},
	}
	err := survey.AskOne(prompt, &state)
	utils.ExitIfError(err)
	mapParams := make(map[string]string)
	mapParams["state"] = state
	jsonParams, err := json.Marshal(mapParams)
	if err != nil {
		log.Fatal("Error ", err)
	}
	responseLogs := request.AirflowPatchRequest("dags/"+DagId+"/dagRuns/"+DagRunId, string(jsonParams))
	var logs model.Logs
	if err := json.Unmarshal([]byte(responseLogs), &logs); err != nil {
		panic(err)
	}
}

func status(cmd *cobra.Command, args []string) {
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

	run, err := prompt.PromptDagRun(DagId, OrderBy, Limit)
	if err != nil {
		panic(err)
	}
	if reflect.DeepEqual(run, model.Dag{}) {
		os.Exit(0)
	}
	DagRunId = run.DagRunId
	dagRunSate := run.State

	var todo string
	promptTodo := &survey.Select{
		Message: "What do you want to do?",
		Options: []string{"View logs", "Change dag run state", "Change task state", "Exit"},
	}

	err = survey.AskOne(promptTodo, &todo)
	utils.ExitIfError(err)

	switch todo {
	case "Exit":
		os.Exit(0)
	case "View logs":
		showLog()
	case "Change dag run state":
		setDagRunState(dagRunSate)
	case "Change task state":
		prompt.PromptSetTaskInstanceState(DagId, DagRunId)
	}
}
