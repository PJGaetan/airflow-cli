package state

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/pkg/model"
	"github.com/pjgaetan/airflow-cli/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	DagId   string
	OrderBy string
)

// NewState represents the list command.
func NewState() *cobra.Command {
	cmd := cobra.Command{
		Use:   "state",
		Short: "Pause/unpause a dag",
		Run:   cmd,
	}
	cmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	cmd.Flags().StringVarP(&OrderBy, "order-by", "o", "-start_date", "The name of the field to order the results by. Prefix a field name with - to reverse the sort order.")
	return &cmd
}

func cmd(cmd *cobra.Command, args []string) {
	if DagId == "" {
		response := request.AirflowGetRequest("dags")
		var dags model.Dags
		if err := json.Unmarshal([]byte(response), &dags); err != nil {
			panic(err)
		}

		if len(dags.Dags) == 0 {
			utils.Failed("No dag found.")
		}

		dagIds := make([]string, len(dags.Dags))
		for index, dag := range dags.Dags {
			dagIds[index] = dag.DagId
		}

		prompt := &survey.Select{
			Message: "Choose a dagRun:",
			Options: dagIds,
			Description: func(value string, index int) string {
				for _, d := range dags.Dags {
					if d.DagId == value {
						if d.IsPaused {
							return "paused"
						}
						return "not paused"
					}
				}
				return ""
			},
		}

		err := survey.AskOne(prompt, &DagId)
		utils.ExitIfError(err)
	}
	if DagId == "" {
		os.Exit(0)
	}

	var todo string
	promptTodo := &survey.Select{
		Message: "What status do you want>",
		Options: []string{"paused", "unpaused"},
	}

	err := survey.AskOne(promptTodo, &todo)
	utils.ExitIfError(err)

	if todo == "" {
		os.Exit(0)
	}

	var isPaused bool
	if todo == "paused" {
		isPaused = true
	} else {
		isPaused = false
	}
	fmt.Println(isPaused)
	mapParams := make(map[string]any)
	mapParams["is_paused"] = isPaused
	jsonParams, err := json.Marshal(mapParams)
	if err != nil {
		log.Fatal("Error ", err)
	}
	response := request.AirflowPatchRequest("dags/"+DagId, string(jsonParams))
	fmt.Println(string(response))
}
