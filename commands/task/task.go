/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"fmt"

	"github.com/pjgaetan/airflow-cli/commands/task/list"
	"github.com/spf13/cobra"
)

func NewTask() *cobra.Command {
	taskCmd := cobra.Command{
		Use:   "task",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("task called")
		},
	}
	taskCmd.AddCommand(list.NewList())
	return &taskCmd
}
