/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package dag

import (
	"github.com/pjgaetan/airflow-cli/commands/dag/list"
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
	return &dagCmd
}
