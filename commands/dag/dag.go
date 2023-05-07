/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package dag

import (
	"fmt"

	"github.com/pjgaetan/airflow-cli/commands/dag/list"
	"github.com/spf13/cobra"
)

func NewDag() *cobra.Command {
	dagCmd := cobra.Command{
		Use:   "dag",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("dag called")
		},
	}
	dagCmd.AddCommand(list.NewList())
	return &dagCmd
}
