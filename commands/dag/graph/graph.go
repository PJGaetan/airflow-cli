package graph

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/awalterschulze/gographviz"
	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/pkg/model"
	"github.com/pjgaetan/airflow-cli/pkg/prompt"
	"github.com/pjgaetan/airflow-cli/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	DagId   string
	OrderBy string
)

// listCmd represents the list command
func NewGraph() *cobra.Command {
	cmd := cobra.Command{
		Use:   "graph",
		Short: "Graphviz representation of the dag",
		Run:   cmd,
	}
	cmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
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
		DagId = dag.Dag_id
	}
	response := request.AirflowGetRequest("dags/" + DagId + "/tasks")
	var tasks model.Tasks
	if err := json.Unmarshal([]byte(response), &tasks); err != nil {
		panic(err)
	}
	if len(tasks.Task) == 0 {
		utils.Failed("No task in this dag")
	}
	graph := gographviz.NewGraph()
	graph.AddAttr(graph.Name, "label", DagId)
	graph.AddAttr(graph.Name, "labelloc", "t")
	graph.AddAttr(graph.Name, "rankdir", "LR")
	// [label=tutorial_dag labelloc=t rankdir=LR]
	graph.SetName(DagId)
	graph.SetDir(true)
	for _, t := range tasks.Task {
		attr := make(map[string]string)
		attr["color"] = "\"#000000\""
		attr["fillcolor"] = "\"#ffefeb\""
		attr["label"] = t.TaskId
		attr["shape"] = "rectangle"
		attr["style"] = "\"filled,rounded\""

		graph.AddNode(DagId, t.TaskId, attr)
		for _, subTask := range t.DownstreamTaskIds {
			graph.AddEdge(t.TaskId, subTask, true, nil)
		}
	}

	output := graph.String()
	fmt.Println(output)
}
