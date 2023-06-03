package prompt

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/pkg/model"
	"github.com/pjgaetan/airflow-cli/pkg/utils"
)

func PromptDag() (model.Dag, error) {
	response := request.AirflowGetRequest("dags")
	var dags model.Dags
	if err := json.Unmarshal([]byte(response), &dags); err != nil {
		panic(err)
	}
	dagIds := make([]string, len(dags.Dags))
	for index, dag := range dags.Dags {
		dagIds[index] = dag.DagId
	}
	if len(dags.Dags) == 0 {
		utils.Failed("No run for this dag.")
	}
	var dagId string
	prompt := &survey.Select{
		Message: "Choose a dagRun:",
		Options: dagIds,
	}

	err := survey.AskOne(prompt, &dagId)
	utils.ExitIfError(err)
	if dagId == "" {
		return model.Dag{}, nil
	}
	for _, dag := range dags.Dags {
		if dag.DagId == dagId {
			return dag, nil
		}
	}
	return model.Dag{}, errors.New("No such a dagId error")
}

func PromptDagRun(dagId, orderBy string, limit int) (model.DagRun, error) {
	response := request.AirflowGetRequest("dags/"+dagId+"/dagRuns", [2]string{"limit", strconv.Itoa(limit)}, [2]string{"order_by", orderBy})
	var dagRun model.DagRuns
	if err := json.Unmarshal([]byte(response), &dagRun); err != nil {
		panic(err)
	}
	runs := dagRun.DagRun

	// reverse sort to have last run in first
	sort.Slice(runs, func(i, j int) bool {
		return runs[i].StartDate.Format(time.RFC3339) > runs[j].StartDate.Format(time.RFC3339)
	})

	dagRunIds := make([]string, len(runs))
	for index, run := range runs {
		dagRunIds[index] = run.DagRunId
	}
	if len(runs) == 0 {
		utils.Failed("No run for this dag.")
	}
	fmt.Println(runs[0].DagId)

	var dagRunId string
	prompt := &survey.Select{
		Message: "Choose a dagRun:",
		Options: dagRunIds,
		Description: func(value string, index int) string {
			for _, run := range runs {
				if run.DagRunId == value {
					if run.State == "" {
						return "none"
					}
					return run.State
				}
			}
			return ""
		},
	}

	err := survey.AskOne(prompt, &dagRunId)
	utils.ExitIfError(err)

	if dagRunId == "" {
		return model.DagRun{}, nil
	}

	for _, run := range runs {
		if run.DagRunId == dagRunId {
			return run, nil
		}
	}
	return model.DagRun{}, errors.New("No such a dagRunId error")
}
