package commands

type Command string

func NewCommand(cmd string) Command {
	return Command(cmd)
}
