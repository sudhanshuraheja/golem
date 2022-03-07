package commands

type Command struct {
	Exec string
	Apt  []Apt
}

type Apt struct {
	PGP              string
	Repository       APTRepository
	Update           bool
	Purge            []string
	Install          []string
	InstallNoUpgrade []string
}

type APTRepository struct {
	URL     string
	Sources string
}

func NewCommand(exec string) *Command {
	return &Command{Exec: exec}
}
