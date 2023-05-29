/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package dagRrun

import (
	"github.com/pjgaetan/airflow-cli/commands/dagRun/list"
	"github.com/pjgaetan/airflow-cli/commands/dagRun/status"
	"github.com/spf13/cobra"
)

func NewRun() *cobra.Command {
	runCmd := cobra.Command{
		Use:   "dag-run",
		Short: "Interact with dag runs",
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println("run called")
		// },
	}
	runCmd.AddCommand(list.NewList())
	runCmd.AddCommand(status.NewStatus())
	return &runCmd
}
