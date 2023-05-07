/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/pkg/model"

	"github.com/spf13/cobra"
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
	return &listCmd
}

func list(cmd *cobra.Command, args []string) {
	dat := request.AirflowGetRequest("dags")
	t := buildTable(dat)

	t.Render()
}

func buildTable(dat model.Dags) table.Writer {
	t := table.NewWriter()
	t.SetStyle(table.Style{
		Name: "pureStyle",
		Box: table.BoxStyle{
			PaddingLeft:  " ",
			PaddingRight: " ",
		},
		Options: table.Options{
			DrawBorder:      false,
			SeparateColumns: false,
			SeparateFooter:  false,
			SeparateHeader:  false,
			SeparateRows:    false,
		},
	})
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
