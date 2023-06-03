package listInstance

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/internal/printer"
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
func NewListInstance() *cobra.Command {
	listCmd := cobra.Command{
		Use:   "list-instance",
		Short: "List task instances linked to a particular dag run",
		Run:   list,
	}
	listCmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	listCmd.Flags().StringVarP(&DagRunId, "dag-run-id", "r", "", "dag id")
	listCmd.Flags().StringVarP(&OrderBy, "order-by", "o", "-start_date", "The name of the field to order the results by. Prefix a field name with - to reverse the sort order.")
	listCmd.Flags().IntVarP(&Limit, "limit", "l", 20, "The numbers of items to return.")
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

	t := buildTable(instances)

	t.Render()
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
