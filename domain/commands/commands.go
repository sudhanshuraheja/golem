package commands

type Commands []Command

func NewCommands(cmds []string) *Commands {
	c := Commands{}
	for _, cmd := range cmds {
		c.Append(NewCommand(cmd))
	}
	return &c
}

func (c *Commands) Append(cmd Command) {
	*c = append(*c, cmd)
}

func (c *Commands) Merge(cmds Commands) {
	*c = append(*c, cmds...)
}

func (c *Commands) ToArray() []Command {
	_cmds := []Command{}
	if c != nil {
		for _, cm := range *c {
			_cmds = append(_cmds, cm)
		}
	}
	return _cmds
}
