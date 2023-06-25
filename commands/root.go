package commands

import (
	"log"

	"github.com/gookit/ini/v2"
	"github.com/spf13/cobra"

	"github.com/pjgaetan/airflow-cli/commands/dag"
	"github.com/pjgaetan/airflow-cli/commands/profile"
	"github.com/pjgaetan/airflow-cli/commands/task"
	"github.com/pjgaetan/airflow-cli/commands/version"
	"github.com/pjgaetan/airflow-cli/internal/config"
	"github.com/pjgaetan/airflow-cli/internal/flag"
	"github.com/pjgaetan/airflow-cli/pkg/utils"
)

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "airflow-cli",
		Short: "Airflow CLI help interract with airflow REST API.",
		Long:  `Airflow CLI interacting with the REST API that helps you manage your dags, tasks.. from the command line.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setConfig(cmd)
		},
	}
	rootCmd.PersistentFlags().StringVarP(&flag.Flag, "profile", "p", "default", "Profile used to connect to Airflow")

	rootCmd.AddCommand(
		profile.NewProfile(),
		dag.NewDag(),
		task.NewTask(),
		version.NewCmdVersion(),
	)
	return &rootCmd
}

func setConfig(cmd *cobra.Command) {
	subCmd := cmd.Name()
	if !cmdRequireToken(subCmd) {
		return
	}
	profile_name, auth_method, err := config.GetActiveProfile()
	if err != nil {
		log.Fatal("Error ", err)
	}
	switch auth_method {
	case "user/password":
		profile := config.GetUserPasswordProfile(profile_name)
		config.AuthorizationHeader = "Authorization: Basic " + config.BasicAuth(profile)
	case "jwt":
		profile := config.GetJwtProfile(profile_name)
		token := config.GetToken(profile)
		config.AuthorizationHeader = "Authorization: Bearer " + token
	default:
		utils.Failed("no such possibility")
	}
	config.Url = ini.String(profile_name + ".url")
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
