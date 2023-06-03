package commands

import (
	"github.com/spf13/cobra"

	"github.com/pjgaetan/airflow-cli/commands/dag"
	"github.com/pjgaetan/airflow-cli/commands/profile"
	"github.com/pjgaetan/airflow-cli/commands/task"
	"github.com/pjgaetan/airflow-cli/internal/flag"
)

var Profile string

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "airflow-cli",
		Short: "Airflow CLI help interract with airflow REST API.",
		Long:  `Airflow CLI interacting with the REST API that helps you manage your dags, tasks.. from the command line.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			checkConfig(cmd)
		},
	}
	rootCmd.PersistentFlags().StringVarP(&flag.Flag, "profile", "p", "default", "Profile used to connect to Airflow")

	rootCmd.AddCommand(
		profile.NewProfile(),
		dag.NewDag(),
		task.NewTask(),
	)
	return &rootCmd
}

func checkConfig(cmd *cobra.Command) {
	subCmd := cmd.Name()
	if !cmdRequireToken(subCmd) {
		return
	}
}

func cmdRequireToken(cmd string) bool {
	allowList := []string{
		"init",
		"help",
		"version",
		"completion",
	}

	for _, item := range allowList {
		if item == cmd {
			return false
		}
	}

	return true
}
