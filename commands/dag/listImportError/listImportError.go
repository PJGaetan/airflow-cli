package listImportError

import (
	"fmt"
	"strconv"

	"github.com/pjgaetan/airflow-cli/api/request"
	"github.com/pjgaetan/airflow-cli/internal/constant"

	"github.com/spf13/cobra"
)

var (
	Limit   int
	OrderBy string
)

func NewListImportError() *cobra.Command {
	cmd := cobra.Command{
		Use:   "list-import-errors",
		Short: "List dags import errors",
		Run:   cmd,
	}
	cmd.Flags().IntVarP(&Limit, "limit", "l", constant.DEAULT_ITEM_LIMIT, "The numbers of items to return.")
	cmd.Flags().StringVarP(&OrderBy, "order-by", "o", "-timestamp", "The name of the field to order the results by. Prefix a field name with - to reverse the sort order.")
	return &cmd
}

func cmd(cmd *cobra.Command, args []string) {
	response := request.AirflowGetRequest("importErrors", [2]string{"limit", strconv.Itoa(Limit)}, [2]string{"order_by", OrderBy})
	fmt.Println(string(response))
}
