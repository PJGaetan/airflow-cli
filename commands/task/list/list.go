package list

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/internal/printer"
	"github.com/pjgaetan/airflow-cli/pkg/model"
	"github.com/pjgaetan/airflow-cli/pkg/prompt"
	"github.com/pjgaetan/airflow-cli/pkg/utils"

	"github.com/spf13/cobra"
)

var DagId string

// listCmd represents the list command.
func NewList() *cobra.Command {
	listCmd := cobra.Command{
		Use:   "list",
		Short: "List task of a specific dag",
		Run:   list,
	}
	listCmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	return &listCmd
}

func list(cmd *cobra.Command, args []string) {
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

	response := request.AirflowGetRequest("dags/" + DagId + "/tasks")
	var tasks model.Tasks
	if err := json.Unmarshal([]byte(response), &tasks); err != nil {
		panic(err)
	}
	t := buildTable(tasks)

	t.Render()
}

func buildTable(dat model.Tasks) table.Writer {
	t := table.NewWriter()
	t = printer.InitTable(t)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"task_id",
		"operator",
		"start_date",
		"end_date",
		"trigger_rule",
		"retry",
		"dowstream_task",
	})
	t.AppendSeparator()
	for _, s := range dat.Task {
		t.AppendRow([]interface{}{
			s.TaskId,
			s.ClassRef.ClassName,
			utils.FormatDate(s.StartDate),
			utils.FormatDate(s.EndDate),
			s.TriggerRule,
			s.Retries,
			"[" + strings.Join(s.DownstreamTaskIds, ",") + "]",
		})
	}
	return t
}
