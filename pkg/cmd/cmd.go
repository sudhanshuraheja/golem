package cmd

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
)

type Cmd struct {
	name   string
	cmd    *exec.Cmd
	stdin  io.Writer
	stdout io.Reader
	stderr io.Reader
	Stdout chan Out
	Stderr chan Out
}

type Out struct {
	Name      string
	Message   string
	Completed bool
}

func NewCmd(name string) *Cmd {
	cmd := Cmd{}
	cmd.name = name
	cmd.Stdout = make(chan Out, 1000)
	cmd.Stderr = make(chan Out, 1000)
	return &cmd
}

func (c *Cmd) Run(command string) error {
	c.cmd = exec.Command("bash", "-c", command)
	err := c.pipes()
	if err != nil {
		return err
	}

	err = c.cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cmd) pipes() error {
	var err error
	c.stdin, err = c.cmd.StdinPipe()
	if err != nil {
		return err
	}

	c.stdout, err = c.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	c.stderr, err = c.cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func(name string) {
		scanner := bufio.NewScanner(c.stdout)
		for {
			if tkn := scanner.Scan(); tkn {
				received := scanner.Bytes()
				data := make([]byte, len(received))
				copy(data, received)
				c.Stdout <- Out{
					Name:    name,
					Message: string(data),
				}
			} else {
				if scanner.Err() != nil {
					errMessage := scanner.Err().Error()
					if strings.Contains(errMessage, "file already closed") {
						c.Stdout <- Out{
							Name:      name,
							Message:   "Done",
							Completed: true,
						}
					} else {
						c.Stderr <- Out{
							Name:      name,
							Message:   scanner.Err().Error(),
							Completed: true,
						}
					}
				} else {
					c.Stdout <- Out{
						Name:      name,
						Message:   "Done",
						Completed: true,
					}
				}
				return
			}
		}
	}(c.name)

	go func(name string) {
		scanner := bufio.NewScanner(c.stderr)
		for scanner.Scan() {
			c.Stderr <- Out{
				Name:    name,
				Message: scanner.Text(),
			}
		}
	}(c.name)

	return nil
}
