/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package dag

import (
	"github.com/pjgaetan/airflow-cli/commands/dag/graph"
	"github.com/pjgaetan/airflow-cli/commands/dag/list"
	"github.com/pjgaetan/airflow-cli/commands/dag/listImportError"
	"github.com/pjgaetan/airflow-cli/commands/dag/listRuns"
	"github.com/pjgaetan/airflow-cli/commands/dag/state"
	"github.com/pjgaetan/airflow-cli/commands/dag/status"
	"github.com/pjgaetan/airflow-cli/commands/dag/trigger"
	"github.com/spf13/cobra"
)

func NewDag() *cobra.Command {
	dagCmd := cobra.Command{
		Use:   "dag",
		Short: "Interact with Dag resources",
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println("dag called")
		// },
	}
	dagCmd.AddCommand(list.NewList())
	dagCmd.AddCommand(trigger.NewTrigger())
	dagCmd.AddCommand(listRuns.NewListRuns())
	dagCmd.AddCommand(status.NewStatus())
	dagCmd.AddCommand(listImportError.NewListImportError())
	dagCmd.AddCommand(state.NewState())
	dagCmd.AddCommand(graph.NewGraph())
	return &dagCmd
}
