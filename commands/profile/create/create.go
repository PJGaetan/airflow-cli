/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/ini/v2"
	"github.com/pjgaetan/airflow-cli/internal/config"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
func NewCreate() *cobra.Command {
	createCmd := cobra.Command{
		Use:   "create",
		Short: "A brief description of your command",
		Long:  ``,
		Run:   create,
	}
	return &createCmd
}

var qJwt = []*survey.Question{
	{
		Name:     "isShell",
		Prompt:   &survey.Confirm{Message: "Is your token a shell command to access it ?"},
		Validate: survey.Required,
	},
	{
		Name:     "token",
		Prompt:   &survey.Input{Message: "Token or command to access the token."},
		Validate: survey.Required,
	},
}

var qUserPassword = []*survey.Question{
	{
		Name:     "user",
		Prompt:   &survey.Input{Message: "Airflow user"},
		Validate: survey.Required,
	},
	{
		Name:     "password",
		Prompt:   &survey.Password{Message: "Airflow password"},
		Validate: survey.Required,
	},
}

func create(cmd *cobra.Command, args []string) {
	profile_name := ""
	prompt := &survey.Input{
		Message: "Profile name",
		Default: "default",
	}
	survey.AskOne(prompt, &profile_name)

	profiles := config.GetProfiles()
	for k := range profiles {
		if k == profile_name {
			override := false
			prompt := &survey.Confirm{
				Message: "Profile " + profile_name + " already exists, do you want to override it ?",
			}
			survey.AskOne(prompt, &override)
			if !override {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("\u001B[0;31m✗\u001B[0m %s\n", "Choosing not to override.\nRerun command to specify another profile name."))
				os.Exit(1)
			}
		}
	}

	url := ""
	prompt_url := &survey.Input{
		Message: "Url to airflow REST API.",
	}
	survey.AskOne(prompt_url, &url)

	auth_type := ""
	auth_prompt := &survey.Select{
		Message: "Chose a way to authentication",
		Options: []string{"user/password", "jwt token"},
		Default: "user/password",
	}
	survey.AskOne(auth_prompt, &auth_type)

	mapProfile := make(map[string]string)
	if auth_type == "user/password" {
		answers := struct {
			User     string
			Password string
		}{}

		err := survey.Ask(qUserPassword, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		mapProfile["url"] = url
		mapProfile["user"] = answers.User
		mapProfile["password"] = answers.Password
	} else {
		answers := struct {
			IsShell bool
			Token   string
		}{}

		err := survey.Ask(qJwt, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		mapProfile["url"] = url
		mapProfile["token"] = answers.Token
		mapProfile["is_shell"] = strconv.FormatBool(answers.IsShell)
	}

	if len(profiles) == 0 {
		fmt.Println("No config file yet.")
	}

	ini.Default().NewSection(profile_name, mapProfile)
	config.WriteConfig()
}
