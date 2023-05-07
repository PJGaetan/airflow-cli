/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"fmt"

	"github.com/gookit/ini/v2"
	"github.com/pjgaetan/airflow-cli/internal/config"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
func NewList() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "List create profiles",
		Long:  ``,
		Run:   profile,
	}
	return &cmd
}

func profile(cmd *cobra.Command, args []string) {
	err := ini.LoadExists("testdata/test.ini", "not-exist.ini")
	if err != nil {
		panic(err)
	}
	fmt.Println(config.GetProfiles())
}
