package printer

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func InitTable(t table.Writer) table.Writer {
	t.SetStyle(table.Style{
		Name: "pureStyle",
		Box: table.BoxStyle{
			PaddingLeft:  " ",
			PaddingRight: " ",
		},
		Options: table.Options{
			DrawBorder:      false,
			SeparateColumns: false,
			SeparateFooter:  false,
			SeparateHeader:  false,
			SeparateRows:    false,
		},
	})
	t.SetOutputMirror(os.Stdout)
	return t
}
