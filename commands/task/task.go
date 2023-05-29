/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"github.com/pjgaetan/airflow-cli/commands/task/list"
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
	return &taskCmd
}
