package log

import (
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type Table struct {
	headerFmt func(format string, a ...interface{}) string
	columnFmt func(format string, a ...interface{}) string
	table     table.Table
}

func NewTable(columns ...interface{}) *Table {
	t := Table{}
	t.headerFmt = color.New(color.FgGreen, color.Underline).SprintfFunc()
	t.columnFmt = color.New(color.FgYellow).SprintfFunc()
	t.table = table.New(columns...)
	t.table.WithHeaderFormatter(t.headerFmt).WithFirstColumnFormatter(t.columnFmt)
	return &t
}

func (t *Table) Row(vals ...interface{}) {
	t.table.AddRow(vals...)
}

func (t *Table) Display() {
	t.table.Print()
}
