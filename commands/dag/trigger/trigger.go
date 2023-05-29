package trigger

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/pkg/model"
	"github.com/pjgaetan/airflow-cli/pkg/prompt"

	"github.com/spf13/cobra"
)

var DagId string

// triggerCmd represents the list command
func NewTrigger() *cobra.Command {
	triggerCmd := cobra.Command{
		Use:   "trigger",
		Short: "Trigger dag runs",
		Run:   trigger,
	}
	triggerCmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	return &triggerCmd
}

func trigger(cmd *cobra.Command, args []string) {
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
	response := request.AirflowPostRequest("dags/"+DagId+"/dagRuns", `{}`)
	var dagRun model.DagRuns
	if err := json.Unmarshal([]byte(response), &dagRun); err != nil {
		panic(err)
	}
	fmt.Println(string(response))
}
