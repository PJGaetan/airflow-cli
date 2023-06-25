package setInstance

import (
	"os"
	"reflect"

	"github.com/pjgaetan/airflow-cli/internal/constant"
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

// NewSet represents the set-instance command.
func NewSet() *cobra.Command {
	stateCmd := cobra.Command{
		Use:   "set-instance",
		Short: "Set task instances state (success,failed)",
		Run:   state,
	}
	stateCmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	stateCmd.Flags().StringVarP(&DagRunId, "dag-run-id", "r", "", "dag id")
	stateCmd.Flags().StringVarP(&OrderBy, "order-by", "o", "-start_date", "The name of the field to order the results by. Prefix a field name with - to reverse the sort order.")
	stateCmd.Flags().IntVarP(&Limit, "limit", "l", constant.DEAULT_ITEM_LIMIT, "The numbers of items to return.")
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
	prompt.PromptSetTaskInstanceState(DagId, DagRunId)
}
