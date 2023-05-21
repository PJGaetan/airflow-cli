/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package dagRrun

import (
	"fmt"

	"github.com/pjgaetan/airflow-cli/commands/dagRun/list"
	"github.com/pjgaetan/airflow-cli/commands/dagRun/status"
	"github.com/spf13/cobra"
)

func NewRun() *cobra.Command {
	runCmd := cobra.Command{
		Use:   "dag-run",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("run called")
		},
	}
	runCmd.AddCommand(list.NewList())
	runCmd.AddCommand(status.NewStatus())
	return &runCmd
}
