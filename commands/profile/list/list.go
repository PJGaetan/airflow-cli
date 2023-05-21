/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"os"

	"github.com/gookit/ini/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pjgaetan/airflow-cli/internal/config"
	"github.com/pjgaetan/airflow-cli/internal/printer"
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
	t := buildTable(config.GetProfiles())
	t.Render()

	// fmt.Println(config.GetProfiles())
}

func buildTable(profiles map[string]ini.Section) table.Writer {
	t := table.NewWriter()
	t = printer.InitTable(t)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"type",
		"name",
		"endpoint",
		"user",
		"isShellToken",
	})
	t.AppendSeparator()
	for name, profile := range profiles {
		var typeProfile string
		if _, ok := profile["user"]; ok {
			typeProfile = "user/password"
		} else {
			typeProfile = "jwt"
		}
		t.AppendRow([]interface{}{
			typeProfile,
			name,
			profile["url"],
			profile["user"],
			profile["isShell"],
		})
	}
	return t
}
