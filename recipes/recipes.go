package recipes

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

func Errors() *color.Color {
	return color.New(color.FgRed, color.Bold)
}

func Success() *color.Color {
	return color.New(color.FgGreen, color.Bold)
}

func Info() *color.Color {
	return color.New(color.FgCyan, color.Bold)
}

func Progress() *color.Color {
	return color.New(color.FgWhite)
}
