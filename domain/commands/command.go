package commands

import (
	"github.com/sudhanshuraheja/golem/domain/template"
)

type Command string

func NewCommand(cmd string) Command {
	return Command(cmd)
}

func (c *Command) PrepareForExecution(tpl *template.Template) (*Command, error) {
	templatePath, err := tpl.Execute(string(*c))
	if err != nil {
		return nil, err
	}
	path := NewCommand(templatePath)
	return &path, nil
}
