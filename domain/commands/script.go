package commands

type Script struct {
	Apt      []Apt      `hcl:"apt,block"`
	Commands *[]Command `hcl:"commands"`
	Command  *Command   `hcl:"command"`
}
