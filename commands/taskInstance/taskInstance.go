package taskInstance

import (
	"github.com/pjgaetan/airflow-cli/commands/taskInstance/clear"
	"github.com/pjgaetan/airflow-cli/commands/taskInstance/list"
	"github.com/pjgaetan/airflow-cli/commands/taskInstance/logs"
	"github.com/pjgaetan/airflow-cli/commands/taskInstance/set"
	"github.com/spf13/cobra"
)

func NewRun() *cobra.Command {
	runCmd := cobra.Command{
		Use:   "task-instance",
		Short: "Interact with task instances",
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println("run called")
		// },
	}
	runCmd.AddCommand(list.NewList())
	runCmd.AddCommand(clear.NewClear())
	runCmd.AddCommand(set.NewSet())
	runCmd.AddCommand(logs.NewLogs())
	return &runCmd
}
