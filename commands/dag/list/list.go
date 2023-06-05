/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/internal/printer"
	"github.com/pjgaetan/airflow-cli/pkg/model"

	"github.com/spf13/cobra"
)

// listCmd represents the list command.
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
		"schedule_interval",
		// "next_dagrun",
		"is_paused",
		"tags",
		"owners",
	})
	t.AppendSeparator()
	for _, s := range dat.Dags {

		tags := make([]string, len(s.Tags))
		for index, t := range s.Tags {
			tags[index] = t.Tag
		}
		tagsFormated := "[" + strings.Join(tags, ",") + "]"

		t.AppendRow([]interface{}{
			s.DagId,
			s.ScheduleInterval.Value,
			// utils.FormatDate(s.NextDagrun),
			s.IsPaused,
			tagsFormated,
			s.Owners,
		})
	}
	return t
}
