/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package status

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/internal/printer"
	"github.com/pjgaetan/airflow-cli/pkg/model"

	"github.com/spf13/cobra"
)

var (
	DagId   string
	Limit   int
	OrderBy string
)

// listCmd represents the list command
func NewStatus() *cobra.Command {
	statusCmd := cobra.Command{
		Use:   "status",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: list,
	}
	statusCmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	statusCmd.Flags().IntVarP(&Limit, "limit", "l", 5, "The numbers of items to return. (Default:5).")
	statusCmd.Flags().StringVarP(&OrderBy, "order-by", "o", "-start_date", "The name of the field to order the results by. Prefix a field name with - to reverse the sort order.")
	return &statusCmd
}

func list(cmd *cobra.Command, args []string) {
	response := request.AirflowGetRequest("dags/"+DagId+"/dagRuns", [2]string{"limit", strconv.Itoa(Limit)}, [2]string{"order_by", OrderBy})
	var dagRun model.DagRuns
	if err := json.Unmarshal([]byte(response), &dagRun); err != nil {
		panic(err)
	}
	runs := dagRun.DagRun

	// reverse sort to have last run in first
	sort.Slice(runs, func(i, j int) bool {
		return runs[i].Start_date.Format(time.RFC3339) > runs[j].Start_date.Format(time.RFC3339)
	})

	dagRunIds := make([]string, len(runs))
	for index, run := range runs {
		dagRunIds[index] = run.Dag_run_id
	}
	if len(runs) == 0 {
		fmt.Println("No run for this dag.")
	}
	fmt.Println(runs[0].Dag_id)

	var dagRunId string
	prompt := &survey.Select{
		Message: "Choose a dagRun:",
		Options: dagRunIds,
		Description: func(value string, index int) string {
			for _, run := range runs {
				return run.State
			}
			return ""
		},
	}

	survey.AskOne(prompt, &dagRunId)

	responseTask := request.AirflowGetRequest("dags/" + DagId + "/dagRuns/" + dagRunId + "/taskInstances")
	var taskInstance model.TaskInstances
	if err := json.Unmarshal([]byte(responseTask), &taskInstance); err != nil {
		panic(err)
	}

	instances := taskInstance.TaskInstance
	// reverse sort to have last run in first
	sort.Slice(instances, func(i, j int) bool {
		return instances[i].StartDate.Format(time.RFC3339) > instances[j].StartDate.Format(time.RFC3339)
	})

	t := buildTable(instances)

	t.Render()

	var todo string
	promptTodo := &survey.Select{
		Message: "What do you want to do?",
		Options: []string{"View logs", "Exit"},
	}

	survey.AskOne(promptTodo, &todo)
	if todo == "Exit" {
		return
	}

	taskInstanceIds := make([]string, len(instances))
	for index, instance := range instances {
		taskInstanceIds[index] = instance.TaskId
	}

	var taskInstanceId string
	promptLogs := &survey.Select{
		Message: "Choose a dagRun:",
		Options: taskInstanceIds,
		Description: func(value string, index int) string {
			for _, instance := range instances {
				return instance.State + ", tryNumber - " + strconv.Itoa(int(instance.TryNumber))
			}
			return ""
		},
	}
	survey.AskOne(promptLogs, &taskInstanceId)
	responseLogs := request.AirflowGetRequest("dags/"+DagId+"/dagRuns/"+dagRunId+"/taskInstances/"+taskInstanceId+"/logs/"+"1", [2]string{"full_content", "true"})
	var logs model.Logs
	if err := json.Unmarshal([]byte(responseLogs), &logs); err != nil {
		panic(err)
	}

	log := strings.Trim(string(logs.Content), "[]()\"")
	for _, line := range strings.Split(log, "\\n") {
		fmt.Println(strings.Trim(line, "\\"))
	}
}

func buildTable(dat []model.TaskInstance) table.Writer {
	t := table.NewWriter()
	t = printer.InitTable(t)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"task_id",
		"state",
		"start_date",
		"end_date",
		"execution_date",
		"queue_when",
	})
	t.AppendSeparator()
	for _, s := range dat {
		t.AppendRow([]interface{}{
			s.TaskId,
			s.State,
			s.StartDate.Format(time.RFC3339),
			s.EndDate.Format(time.RFC3339),
			s.ExecutionDate.Format(time.RFC3339),
			s.QueuedWhen,
		})
	}
	return t
}
