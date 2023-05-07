/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package profile

import (
	"github.com/pjgaetan/airflow-cli/commands/profile/create"
	"github.com/pjgaetan/airflow-cli/commands/profile/list"
	"github.com/spf13/cobra"
)

func NewProfile() *cobra.Command {
	profileCmd := cobra.Command{
		Use:   "profile <command> [flags]",
		Short: "Create, list, switch airflow profiles.",
		Long:  ``,
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println("profile called")
		// },
	}
	profileCmd.AddCommand(create.NewCreate(), list.NewList())
	return &profileCmd
}
