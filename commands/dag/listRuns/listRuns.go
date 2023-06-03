package listRuns

import (
	"encoding/json"
	"os"
	"reflect"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/internal/constant"
	"github.com/pjgaetan/airflow-cli/internal/printer"
	"github.com/pjgaetan/airflow-cli/pkg/model"
	"github.com/pjgaetan/airflow-cli/pkg/prompt"
	"github.com/pjgaetan/airflow-cli/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	DagId   string
	Limit   int
	OrderBy string
)

// NewListRuns represents the list command.
func NewListRuns() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list-runs",
		Short: "List dag runs",
		Run:   cmd,
	}
	cmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	cmd.Flags().IntVarP(&Limit, "limit", "l", constant.DEAULT_ITEM_LIMIT, "The numbers of items to return.")
	cmd.Flags().StringVarP(&OrderBy, "order-by", "o", "-start_date", "The name of the field to order the results by. Prefix a field name with - to reverse the sort order.")
	return &cmd
}

func cmd(cmd *cobra.Command, args []string) {
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

	response := request.AirflowGetRequest("dags/"+DagId+"/dagRuns", [2]string{"limit", strconv.Itoa(Limit)}, [2]string{"order_by", OrderBy})
	var dagRun model.DagRuns
	if err := json.Unmarshal([]byte(response), &dagRun); err != nil {
		panic(err)
	}
	t := buildTable(dagRun)

	t.Render()
}

func buildTable(dat model.DagRuns) table.Writer {
	t := table.NewWriter()
	t = printer.InitTable(t)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"dag_id",
		"dag_run_id",
		"extrnale_trigger",
		"state",
		"start_date",
		"end_date",
	})
	t.AppendSeparator()
	for _, s := range dat.DagRun {
		t.AppendRow([]interface{}{
			s.DagId,
			s.DagRunId,
			s.ExternalTrigger,
			s.State,
			utils.FormatDate(s.StartDate),
			utils.FormatDate(s.EndDate),
		})
	}
	return t
}
