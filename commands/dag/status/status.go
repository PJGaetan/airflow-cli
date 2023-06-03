/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package status

import (
	"encoding/json"
	"fmt"
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

	"github.com/spf13/cobra"
)

var (
	DagId    string
	DagRunId string
	Limit    int
	OrderBy  string
)

// listCmd represents the list command
func NewStatus() *cobra.Command {
	statusCmd := cobra.Command{
		Use:   "status",
		Short: "Status of dag runs (tasks & logs)",
		Run:   status,
	}
	statusCmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	statusCmd.Flags().StringVarP(&DagRunId, "dag-run", "r", "", "dag run id")
	statusCmd.Flags().IntVarP(&Limit, "limit", "l", 5, "The numbers of items to return. (Default:5).")
	statusCmd.Flags().StringVarP(&OrderBy, "order-by", "o", "-start_date", "The name of the field to order the results by. Prefix a field name with - to reverse the sort order.")
	return &statusCmd
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

	var todo string
	promptTodo := &survey.Select{
		Message: "What do you want to do?",
		Options: []string{"View logs", "Exit"},
	}

	survey.AskOne(promptTodo, &todo)
	if todo == "Exit" {
		os.Exit(0)
	}

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
	survey.AskOne(promptLogs, &taskInstanceId)
	responseLogs := request.AirflowGetRequest("dags/"+DagId+"/dagRuns/"+DagRunId+"/taskInstances/"+taskInstanceId+"/logs/"+"1", [2]string{"full_content", "true"})
	var logs model.Logs
	if err := json.Unmarshal([]byte(responseLogs), &logs); err != nil {
		panic(err)
	}

	log := strings.Trim(string(logs.Content), "[]()\"")
	for _, line := range strings.Split(log, "\\n") {
		fmt.Println(strings.Trim(line, "\\"))
	}
}
