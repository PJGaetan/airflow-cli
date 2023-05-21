/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"encoding/json"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/internal/printer"
	"github.com/pjgaetan/airflow-cli/pkg/model"

	"github.com/spf13/cobra"
)

var DagId string

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
	return &listCmd
}

func list(cmd *cobra.Command, args []string) {
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
		"dowstream_task",
	})
	t.AppendSeparator()
	for _, s := range dat.Task {
		t.AppendRow([]interface{}{
			s.TaskId,
			s.ClassRef.ClassName,
			s.StartDate.Format(time.RFC3339),
			s.EndDate.Format(time.RFC3339),
			string(s.DownstreamTaskIds),
		})
	}
	return t
}
