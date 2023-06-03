package logs

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
	"github.com/pjgaetan/airflow-cli/internal/constant"
	"github.com/pjgaetan/airflow-cli/pkg/model"
	"github.com/pjgaetan/airflow-cli/pkg/prompt"
	"github.com/pjgaetan/airflow-cli/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	DagId    string
	DagRunId string
	TaskId   string
	Limit    int
	OrderBy  string
)

func NewLogs() *cobra.Command {
	logsCmd := cobra.Command{
		Use:   "logs",
		Short: "Get logs from task instances",
		Run:   logs,
	}
	logsCmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	logsCmd.Flags().StringVarP(&DagRunId, "dag-run-id", "r", "", "dag run id")
	logsCmd.Flags().StringVarP(&DagRunId, "task", "t", "", "task id")
	logsCmd.Flags().StringVarP(&OrderBy, "order-by", "o", "-start_date", "The name of the field to order the results by. Prefix a field name with - to reverse the sort order.")
	logsCmd.Flags().IntVarP(&Limit, "limit", "l", constant.DEAULT_ITEM_LIMIT, "The numbers of items to return.")
	return &logsCmd
}

func logs(cmd *cobra.Command, args []string) {
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

	var instance model.TaskInstance
	for _, i := range instances {
		if i.TaskId == taskInstanceId {
			instance = i
			break
		}
	}
	if reflect.DeepEqual(instance, model.TaskInstance{}) {
		os.Exit(0)
	}
	var logNumber string
	if instance.TryNumber > 1.0 {
		s := make([]string, int(instance.TryNumber))
		for index := range s {
			s[index] = strconv.Itoa(index + 1)
		}
		fmt.Println(s)
		promptLogsNumber := &survey.Select{
			Message: "Which log try?:",
			Options: s,
			Default: '1',
		}
		err := survey.AskOne(promptLogsNumber, &logNumber)
		utils.ExitIfError(err)
	} else {
		logNumber = "1"
	}

	responseLogs := request.AirflowGetRequest("dags/"+DagId+"/dagRuns/"+DagRunId+"/taskInstances/"+instance.TaskId+"/logs/"+logNumber, [2]string{"full_content", "true"})
	var logs model.Logs
	if err := json.Unmarshal([]byte(responseLogs), &logs); err != nil {
		panic(err)
	}

	log := strings.Trim(string(logs.Content), "[]()\"")
	for _, line := range strings.Split(log, "\\n") {
		fmt.Println(strings.Trim(line, "\\"))
	}
}
