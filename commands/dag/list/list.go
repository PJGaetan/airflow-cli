/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"encoding/json"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/internal/printer"
	"github.com/pjgaetan/airflow-cli/pkg/model"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
func NewList() *cobra.Command {
	listCmd := cobra.Command{
		Use:   "list",
		Short: "list dags",
		Run:   list,
	}
	return &listCmd
}

func list(cmd *cobra.Command, args []string) {
	response := request.AirflowGetRequest("dags")
	var dag model.Dags
	if err := json.Unmarshal([]byte(response), &dag); err != nil {
		panic(err)
	}
	t := buildTable(dag)

	t.Render()
}

func buildTable(dat model.Dags) table.Writer {
	t := table.NewWriter()
	t = printer.InitTable(t)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"dag_id",
		// "description",
		// "file_token",
		// "fileloc",
		"is_active",
		"is_paused",
		// "s_subdag",
		"owners",
		// "root_dag_id",
		"schedule_interval",
		// "tags",
	})
	t.AppendSeparator()
	for _, s := range dat.Dags {
		t.AppendRow([]interface{}{
			s.Dag_id,
			// s.Description,
			// s.File_token,
			// s.Fileloc,
			s.Is_active,
			s.Is_paused,
			// s.S_subdag,
			s.Owners,
			// s.Root_dag_id,
			s.Schedule_interval.Value,
			// s.Tags,
		})
	}
	return t
}
