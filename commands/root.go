/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package commands

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pjgaetan/airflow-cli/commands/dag"
	"github.com/pjgaetan/airflow-cli/commands/dagRun"
	"github.com/pjgaetan/airflow-cli/commands/profile"
	"github.com/pjgaetan/airflow-cli/commands/task"
	"github.com/pjgaetan/airflow-cli/internal/flag"
)

var (
	config  string
	debug   bool
	Profile string
)

func init() {
	// cobra.OnInitialize(initConfig)
}

func NewRootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "airflow-cli",
		Short: "Airflow CLI help interract with airflow REST API.",
		Long: `Airflow CLI that helps you manage your dags, tasks.. from the
	command line.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			checkConfig(cmd, args)
		},
	}
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&flag.Flag, "profile", "p", "default", "profile (default: default)")

	rootCmd.AddCommand(
		dag.NewDag(),
		profile.NewProfile(),
		task.NewTask(),
		dagRrun.NewRun(),
	)
	return &rootCmd
}

func checkConfig(cmd *cobra.Command, args []string) {
	subCmd := cmd.Name()
	if !cmdRequireToken(subCmd) {
		return
	}

	// TODO: Put be on when adding config
	// if err := viper.ReadInConfig(); err != nil {
	// 	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	// 		fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;31m✗\u001B[0m %s\n", "Missing configuration file.\nRun 'airflow-cli configure' to configure the tool."))
	// 		os.Exit(1)
	// 	}
	// }
}

func initConfig() {
	if config != "" {
		viper.SetConfigFile(config)
	} else {

		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;31m✗\u001B[0m %s\n", err))
			os.Exit(1)
		}

		viper.AddConfigPath(fmt.Sprintf("%s/%s/%s", home, ".config", ".airflow"))
		viper.SetConfigName(".config")
		viper.SetConfigType("yml")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("airflow_cli")

	if err := viper.ReadInConfig(); err == nil && debug {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
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
