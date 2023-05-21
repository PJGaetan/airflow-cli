/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
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
func NewList() *cobra.Command {
	listCmd := cobra.Command{
		Use:   "list",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: list,
	}
	listCmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	listCmd.Flags().IntVarP(&Limit, "limit", "l", 10, "The numbers of items to return. (Default:10).")
	listCmd.Flags().StringVarP(&OrderBy, "order-by", "o", "-start_date", "The name of the field to order the results by. Prefix a field name with - to reverse the sort order.")
	return &listCmd
}

func list(cmd *cobra.Command, args []string) {
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
		// "conf",
	})
	t.AppendSeparator()
	for _, s := range dat.DagRun {
		t.AppendRow([]interface{}{
			s.Dag_id,
			s.Dag_run_id,
			s.External_trigger,
			s.State,
			s.Start_date.Format(time.RFC3339),
			s.End_date.Format(time.RFC3339),
			// s.Execution_date,
			// s.Logical_date,
			// s.Conf,
		})
	}
	return t
}
