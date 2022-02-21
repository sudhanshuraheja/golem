package ssh

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type Connection struct {
	name        string
	conn        *ssh.Client
	sshSession  *ssh.Session
	sftpSession *sftp.Client
	stdin       io.Writer
	stdout      io.Reader
	stderr      io.Reader
	Stdout      chan Out
	Stderr      chan Out
}

type Out struct {
	Name      string
	Message   string
	Completed bool
}

func NewSSHConnection(name, user, host string, port int, privateKeyPath string) (*Connection, error) {
	connection := Connection{}
	connection.name = name
	connection.Stdout = make(chan Out, 1000)
	connection.Stderr = make(chan Out, 1000)

	err := connection.dial(user, host, port, privateKeyPath)
	if err != nil {
		return nil, err
	}
	return &connection, nil
}

func (c *Connection) dial(user, host string, port int, privateKeyPath string) error {
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
		Timeout: 5 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	c.conn, err = ssh.Dial("tcp", addr, conf)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) getSSHSession() error {
	var err error
	c.sshSession, err = c.conn.NewSession()
	if err != nil {
		return err
	}

	err = c.pty()
	if err != nil {
		return err
	}

	err = c.pipes()
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) pty() error {
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	err := c.sshSession.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) pipes() error {
	var err error
	c.stdin, err = c.sshSession.StdinPipe()
	if err != nil {
		return err
	}

	c.stdout, err = c.sshSession.StdoutPipe()
	if err != nil {
		return err
	}

	c.stderr, err = c.sshSession.StderrPipe()
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
					c.Stderr <- Out{
						Name:      name,
						Message:   scanner.Err().Error(),
						Completed: true,
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

func (c *Connection) Run(command string) (int, error) {
	exitStatus := -1

	err := c.getSSHSession()
	if err != nil {
		return exitStatus, err
	}

	err = c.sshSession.Run(command)
	if err != nil {
		switch v := err.(type) {
		case *ssh.ExitError:
			return v.Waitmsg.ExitStatus(), nil
		default:
			return exitStatus, err
		}
	}

	c.sshSession.Close()
	return exitStatus, err
}

func (c *Connection) getSFTPSession() error {
	var err error
	// c.sftpSession, err = sftp.NewClient(c.conn, sftp.MaxPacket(1e9))
	c.sftpSession, err = sftp.NewClient(c.conn)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) Upload(src, dest string) (int64, error) {
	err := c.getSFTPSession()
	if err != nil {
		return 0, err
	}

	d, err := c.sftpSession.OpenFile(dest, syscall.O_RDWR|syscall.O_CREAT)
	if err != nil {
		return 0, err
	}
	defer d.Close()

	s, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	s_info, err := s.Stat()
	if err != nil {
		return 0, err
	}
	sourceSize := s_info.Size()
	defer s.Close()

	copied, err := io.Copy(d, s)
	if err != nil {
		return 0, err
	}
	if copied != sourceSize {
		return 0, fmt.Errorf("only %d bytes out of %d were copied from %s to %s", copied, sourceSize, src, dest)
	}

	c.sftpSession.Close()
	return copied, nil
}

func (c *Connection) Close() {
	c.conn.Close()
}
