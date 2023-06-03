/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"github.com/pjgaetan/airflow-cli/commands/task/clearInstance"
	"github.com/pjgaetan/airflow-cli/commands/task/list"
	"github.com/pjgaetan/airflow-cli/commands/task/listInstance"
	"github.com/pjgaetan/airflow-cli/commands/task/logs"
	"github.com/pjgaetan/airflow-cli/commands/task/setInstance"
	"github.com/spf13/cobra"
)

func NewTask() *cobra.Command {
	taskCmd := cobra.Command{
		Use:   "task",
		Short: "Interact with tasks",
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println("task called")
		// },
	}
	taskCmd.AddCommand(list.NewList())
	taskCmd.AddCommand(logs.NewLogs())
	taskCmd.AddCommand(listInstance.NewListInstance())
	taskCmd.AddCommand(setInstance.NewSet())
	taskCmd.AddCommand(clearInstance.NewClear())
	return &taskCmd
}
