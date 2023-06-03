package create

import (
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/ini/v2"
	"github.com/pjgaetan/airflow-cli/internal/config"
	"github.com/pjgaetan/airflow-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// NewCreate represents the create command.
func NewCreate() *cobra.Command {
	createCmd := cobra.Command{
		Use:   "create",
		Short: "Create a profile to connect to airflow REST API",
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
	profileName := ""
	prompt := &survey.Input{
		Message: "Profile name",
		Default: "default",
	}

	err := survey.AskOne(prompt, &profileName)
	utils.ExitIfError(err)

	if profileName == "" {
		os.Exit(0)
	}

	profiles := config.GetProfiles()
	for k := range profiles {
		if k == profileName {
			override := false
			prompt := &survey.Confirm{
				Message: "Profile " + profileName + " already exists, do you want to override it ?",
			}
			err := survey.AskOne(prompt, &override)
			utils.ExitIfError(err)
			if !override {
				utils.Failed("Choosing not to override.\nRerun command to specify another profile name.")
			}
		}
	}

	url := ""
	prompt_url := &survey.Input{
		Message: "Url to airflow REST API.",
	}
	err = survey.AskOne(prompt_url, &url)
	utils.ExitIfError(err)

	if url == "" {
		os.Exit(0)
	}

	authType := ""
	authPrompt := &survey.Select{
		Message: "Chose a way to authentication",
		Options: []string{"user/password", "jwt token"},
		Default: "user/password",
	}
	err = survey.AskOne(authPrompt, &authType)
	utils.ExitIfError(err)

	if authType == "" {
		os.Exit(0)
	}

	mapProfile := make(map[string]string)
	if authType == "user/password" {
		answers := struct {
			User     string
			Password string
		}{}

		err := survey.Ask(qUserPassword, &answers)
		utils.ExitIfError(err)

		mapProfile["url"] = url
		mapProfile["user"] = answers.User
		mapProfile["password"] = answers.Password
	} else {
		answers := struct {
			IsShell bool
			Token   string
		}{}

		err := survey.Ask(qJwt, &answers)
		utils.ExitIfError(err)

		mapProfile["url"] = url
		mapProfile["token"] = answers.Token
		mapProfile["is_shell"] = strconv.FormatBool(answers.IsShell)
	}

	if len(profiles) == 0 {
		fmt.Println("No config file yet.")
	}

	errProfile := ini.Default().NewSection(profileName, mapProfile)
	utils.ExitIfError(errProfile)
	config.WriteConfig()
}
