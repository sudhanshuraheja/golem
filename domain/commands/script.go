package commands

type Script struct {
	Apt      []Apt
	Commands Commands
	Command  *Command
}
