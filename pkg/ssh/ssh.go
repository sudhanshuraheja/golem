package ssh

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type Connection struct {
	conn      *ssh.Client
	session   *ssh.Session
	stdin     io.Writer
	stdout    io.Reader
	stderr    io.Reader
	Stdout    chan string
	Stderr    chan string
	Completed chan struct{}
}

func NewSSHConnection(user, host string, port int, privateKeyPath string) (*Connection, error) {
	connection := Connection{}
	connection.Stdout = make(chan string, 1000)
	connection.Stderr = make(chan string, 1000)
	connection.Completed = make(chan struct{})

	err := connection.Dial(user, host, port, privateKeyPath)
	if err != nil {
		return nil, err
	}
	return &connection, nil
}

func (c *Connection) Dial(user, host string, port int, privateKeyPath string) error {
	if user == "" {
		return fmt.Errorf("user cannot be empty")
	}

	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	if port == 0 {
		port = 22
	}

	if privateKeyPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		privateKeyPath = fmt.Sprintf("%s/.ssh/id_rsa", homeDir)
	}

	pemBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	knownHosts := fmt.Sprintf("%s/.ssh/known_hosts", homeDir)

	hostKeyCallback, err := knownhosts.New(knownHosts)
	if err != nil {
		return err
	}

	conf := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: hostKeyCallback,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	c.conn, err = ssh.Dial("tcp", addr, conf)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) Session() error {
	var err error
	c.session, err = c.conn.NewSession()
	if err != nil {
		return err
	}

	err = c.Pty()
	if err != nil {
		return err
	}

	err = c.Pipes()
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) Pty() error {
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	err := c.session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) Pipes() error {
	var err error
	c.stdin, err = c.session.StdinPipe()
	if err != nil {
		return err
	}

	c.stdout, err = c.session.StdoutPipe()
	if err != nil {
		return err
	}

	c.stderr, err = c.session.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(c.stdout)
		for {
			if tkn := scanner.Scan(); tkn {
				received := scanner.Bytes()
				data := make([]byte, len(received))
				copy(data, received)
				c.Stdout <- string(data)
			} else {
				if scanner.Err() != nil {
					c.Stderr <- scanner.Err().Error()
				} else {
					c.Completed <- struct{}{}
				}
				return
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(c.stderr)
		for scanner.Scan() {
			c.Stderr <- scanner.Text()
		}
	}()

	return nil
}

func (c *Connection) Run(command string) (int, error) {
	exitStatus := -1

	err := c.Session()
	if err != nil {
		return exitStatus, err
	}

	err = c.session.Run(command)
	if err != nil {
		switch v := err.(type) {
		case *ssh.ExitError:
			return v.Waitmsg.ExitStatus(), nil
		default:
			return exitStatus, err
		}
	}

	return exitStatus, err
}

func (c *Connection) Close() {
	c.session.Close()
	c.conn.Close()
}
