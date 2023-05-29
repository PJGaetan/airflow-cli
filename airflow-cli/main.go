/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/pjgaetan/airflow-cli/commands"
	"github.com/pjgaetan/airflow-cli/pkg/utils"
)

func main() {
	cmd := commands.NewRootCmd()
	if _, err := cmd.ExecuteC(); err != nil {
		utils.ExitIfError(err)
	}
}
