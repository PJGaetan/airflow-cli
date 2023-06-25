package grid

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"sync"

	"github.com/fatih/color"

	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/internal/constant"
	"github.com/pjgaetan/airflow-cli/pkg/model"
	"github.com/pjgaetan/airflow-cli/pkg/prompt"
	"github.com/pjgaetan/airflow-cli/pkg/utils"
	"golang.org/x/exp/slices"

	"github.com/spf13/cobra"
)

var (
	DagId   string
	Limit   int
	OrderBy string
)

// NewGrid represents the grid command.
func NewGrid() *cobra.Command {
	cmd := cobra.Command{
		Use:   "grid",
		Short: "Dashbord of the results of you lastest dag runs (like the graph view in Airflow UI).",
		Run:   cmd,
	}
	cmd.Flags().StringVarP(&DagId, "dag-id", "d", "", "dag id")
	cmd.Flags().IntVarP(&Limit, "limit", "l", constant.DEAULT_ITEM_LIMIT, "The numbers of items to return.")
	cmd.Flags().StringVarP(&OrderBy, "order-by", "o", "-start_date", "The name of the field to order the results by. Prefix a field name with - to reverse the sort order.")
	return &cmd
}

func asyncGetTaskInstances(sortedTasks []string, dagRun model.DagRuns) [][]string {
	sortedTaskMap := make(map[string]int)
	for index, task := range sortedTasks {
		sortedTaskMap[task] = index
	}

	var chans []chan []model.TaskInstance
	for range dagRun.DagRun {
		chans = append(chans, make(chan []model.TaskInstance))
	}

	var wg sync.WaitGroup
	for index, d := range dagRun.DagRun {
		go getTaskInstances(DagId, d.DagRunId, sortedTaskMap, chans[index])
	}
	wg.Wait()

	a := make([][]bool, len(dagRun.DagRun))
	for i := range a {
		a[i] = make([]bool, len(sortedTasks))
	}
	taskStates := make([][]string, 0)
	for _, t := range chans {
		taskInstances := <-t
		internalSlice := make([]string, 0)

		sortedTaskIndex := 0
		for taskInstanceIndex := 0; taskInstanceIndex < len(sortedTasks); taskInstanceIndex++ {

			// Task was run in past but not present anymore in the dag
			// the task is not present in current version of the code
			// Task does not show (similar with airflow behaviour)
			if !slices.Contains(sortedTasks, taskInstances[taskInstanceIndex].TaskId) {
				continue
			}

			// Missing tasks
			// Happen when new task is added for example
			if taskInstances[taskInstanceIndex].TaskId != sortedTasks[taskInstanceIndex] {
				internalSlice = append(internalSlice, "none")
				sortedTaskIndex++
				continue
			}

			internalSlice = append(internalSlice, taskInstances[taskInstanceIndex].State)
		}
		taskStates = append(taskStates, internalSlice)
	}
	return taskStates
}

func revertMatrix(m [][]string) ([][]string, error) {
	if len(m) == 0 {
		return nil, errors.New("Array to matrix reverse is empty.")
	}

	revertedMatrix := make([][]string, len(m[0]))
	for index := range m[0] {
		revertedMatrix[index] = make([]string, len(m))
	}

	for i := 0; i < len(revertedMatrix); i++ {
		for j := 0; j < len(revertedMatrix[0]); j++ {
			revertedMatrix[i][len(revertedMatrix[0])-j-1] = m[j][i]
		}
	}
	return revertedMatrix, nil
}

func formatState(state string) string {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	cyan := color.New(color.FgCyan)
	switch state {
	case "success":
		return green.Sprintf("v")
	case "running":
		return color.New(color.BgHiGreen).Add(color.FgBlack).Sprintf("r")
	case "failed":
		return red.Sprintf("x")
	case "upstream_failed":
		return color.New(color.FgRed).Add(color.Faint).Sprintf("x")
	case "skipped":
		return yellow.Sprintf("-")
	case "up_for_retry":
		return color.New(color.FgBlack).Add(color.BgHiYellow).Sprintf("r")
	case "up_for_reschedule":
		return color.New(color.FgBlack).Add(color.BgHiYellow).Sprintf("r")
	case "queued":
		return color.New(color.FgHiBlack).Sprintf("q")
	case "none":
		return " "
	case "scheduled":
		return cyan.Sprintf("s")
	case "deferred":
		return " "
	case "removed":
		return " "
	case "restarting":
		return " "
	default:
		return " "
	}
}

// Compute right rightPadding
// take the max len between dag name and tasks names.
func rightPadding(tasks []string, dagName string) int {
	maxTaskLen := 0
	for _, t := range tasks {
		if len(t) > maxTaskLen {
			maxTaskLen = len(t)
		}
	}
	if len(dagName) > maxTaskLen {
		maxTaskLen = len(dagName)
	}
	ADDITIONAL_PADDING := 4
	return maxTaskLen + ADDITIONAL_PADDING
}

/*
dag_id
v x v x v x v x   task1
v x v x v x v x   task2
v x v x v x v x   task3

	|       |

03-12   04-11
*/
func printGraph(states [][]string, orderedTasks []string, rightPadding int) {
	for index, tasks := range states {
		fmt.Printf("%-*s", rightPadding, orderedTasks[index])
		for _, state := range tasks {
			fmt.Print(formatState(state), " ")
		}
		fmt.Println("")
	}
}

func getTaskInstances(dagId string, dagRunId string, sortedTaskMap map[string]int, c chan<- []model.TaskInstance) {
	responseTask := request.AirflowGetRequest("dags/" + dagId + "/dagRuns/" + dagRunId + "/taskInstances")
	var taskInstance model.TaskInstances
	if err := json.Unmarshal([]byte(responseTask), &taskInstance); err != nil {
		panic(err)
	}

	instances := taskInstance.TaskInstance
	sort.Slice(instances, func(i, j int) bool {
		return sortedTaskMap[instances[i].TaskId] < sortedTaskMap[instances[j].TaskId]
	})
	c <- instances
}

func printHeader(dagRuns []model.DagRun, rightPadding int) {
	fmt.Printf("%-*s", rightPadding, dagRuns[0].DagId)
	for i := 0; i < len(dagRuns); i++ {
		run := dagRuns[len(dagRuns)-i-1]
		switch s := run.State; s {
		case "success":
			color.New(color.FgGreen).Add(color.Bold).Print("v", " ")
		case "failed":
			color.New(color.FgRed).Add(color.Bold).Print("x", " ")
		default:
			fmt.Print("  ")
		}
	}
	fmt.Println("")
	for i := 0; i < 2*len(dagRuns)+rightPadding-1; i++ {
		fmt.Print("-")
	}
	fmt.Println("")
}

func printFooter(dagRuns []model.DagRun, rightPadding int) {
	TASK_BETWEEN_DATE := 7
	STEP := 2

	// traverse and put "|" under last run and every TASK_BETWEEN_DATE run before it.
	for i := 0; i < len(dagRuns)+rightPadding; i++ {
		if i < rightPadding {
			fmt.Print(" ")
			continue
		}
		if (i-rightPadding-len(dagRuns)+1) == 0 || (i-rightPadding-len(dagRuns)+1)%TASK_BETWEEN_DATE == 0 {
			fmt.Print("|")
			fmt.Print(" ")
			continue
		}
		fmt.Print("  ")

	}
	fmt.Println("")

	// put date under "|"
	LEN_DATE := 10
	for i := 0; i < len(dagRuns)+rightPadding; i++ {
		if i < rightPadding {
			fmt.Print(" ")
			continue
		}
		dagRunInstanceOrdered := rightPadding + len(dagRuns) - 1 - (LEN_DATE / (2 * STEP)) - i
		if dagRunInstanceOrdered == 0 {
			fmt.Printf(dagRuns[0].ExecutionDate.Format("2006-01-02"))
			continue
		}
		if dagRunInstanceOrdered%TASK_BETWEEN_DATE == 0 {
			umpteenthTask := dagRunInstanceOrdered / TASK_BETWEEN_DATE
			fmt.Printf(dagRuns[TASK_BETWEEN_DATE*umpteenthTask].ExecutionDate.Format("2006-01-02"))
			i += LEN_DATE/2 - 1
			continue
		}
		fmt.Print("  ")

	}
	fmt.Println("")
	// put time under "|"
	LEN_TIME := 8
	for i := 0; i < len(dagRuns)+rightPadding; i++ {
		if i < rightPadding {
			fmt.Print(" ")
			continue
		}
		dagRunInstanceOrdered := rightPadding + len(dagRuns) - 1 - (LEN_TIME / (2 * STEP)) - i
		if dagRunInstanceOrdered == 0 {
			fmt.Printf(dagRuns[0].ExecutionDate.Format("15:04:05"))
			continue
		}
		if dagRunInstanceOrdered%TASK_BETWEEN_DATE == 0 {
			umpteenthTask := dagRunInstanceOrdered / TASK_BETWEEN_DATE
			fmt.Printf(dagRuns[TASK_BETWEEN_DATE*umpteenthTask].ExecutionDate.Format("15:04:05"))
			i += LEN_TIME/2 - 1
			continue
		}
		fmt.Print("  ")

	}
	fmt.Println("")
}

func cmd(cmd *cobra.Command, args []string) {
	if DagId == "" {
		dag, err := prompt.PromptDag()
		utils.ExitIfError(err)
		if reflect.DeepEqual(dag, model.Dag{}) {
			os.Exit(0)
		}
		DagId = dag.DagId
	}

	response := request.AirflowGetRequest("dags/"+DagId+"/dagRuns", [2]string{"limit", strconv.Itoa(Limit)}, [2]string{"order_by", OrderBy})
	var dagRun model.DagRuns
	if err := json.Unmarshal([]byte(response), &dagRun); err != nil {
		panic(err)
	}

	response = request.AirflowGetRequest("dags/" + DagId + "/tasks")
	var tasks model.Tasks
	if err := json.Unmarshal([]byte(response), &tasks); err != nil {
		panic(err)
	}
	nodes := make(map[string]utils.Node)
	for index, t := range tasks.Task {
		nodes[t.TaskId] = utils.Node{
			Id:             int64(index),
			Name:           t.TaskId,
			DownStreamName: t.DownstreamTaskIds,
		}
	}

	sortedTasks := utils.TopoSort(nodes)
	taskStates := asyncGetTaskInstances(sortedTasks, dagRun)

	states, err := revertMatrix(taskStates)
	utils.ExitIfError(err)
	rightPadding := rightPadding(sortedTasks, DagId)
	printHeader(dagRun.DagRun, rightPadding)
	printGraph(states, sortedTasks, rightPadding)
	printFooter(dagRun.DagRun, rightPadding)
}
