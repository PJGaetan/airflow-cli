/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/pjgaetan/airflow-cli/commands"
)

func main() {
	cmd := commands.NewRootCmd()
	if _, err := cmd.ExecuteC(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
