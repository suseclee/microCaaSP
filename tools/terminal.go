package tools

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/suseclee/microCaaSP/configs/constants"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type Node struct {
	user     string
	password string
	port     string
	host     string
}

func (n *Node) Init() {
	n.user = constants.USERNAME
	n.password = constants.GetPassword()
	n.port = "22"
	n.host = constants.NODEIP
}

//https://gist.github.com/atotto/ba19155295d95c8d75881e145c751372
func Terminal() {
	n := &Node{}
	n.Init()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := n.run(ctx); err != nil {
			log.Print(err)
		}
		cancel()
	}()

	select {
	case <-sig:
		cancel()
	case <-ctx.Done():
	}
}

func WaitForLogin() error {
	waitTimes := []int{7, 5, 7, 5, 10, 6}
	waitTime := 0
	for _, duration := range waitTimes {
		ShellSpin([]string{"sleep", strconv.Itoa(duration)})
		waitTime += duration
		if err := Ping(); err == nil {
			return nil
		}
	}
	return errors.Unwrap(fmt.Errorf("Could not boot within %d secs", waitTime))
}

func Ping() error {
	n := &Node{}
	n.Init()

	config := &ssh.ClientConfig{
		User: n.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(n.password),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	hostport := fmt.Sprintf("%s:%s", n.host, n.port)
	conn, err := ssh.Dial("tcp", hostport, config)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func (n *Node) run(ctx context.Context) error {
	config := &ssh.ClientConfig{
		User: n.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(n.password),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	hostport := fmt.Sprintf("%s:%s", n.host, n.port)
	conn, err := ssh.Dial("tcp", hostport, config)
	if err != nil {
		return fmt.Errorf("cannot connect %v: %v", hostport, err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("cannot open new session: %v", err)
	}
	defer session.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("terminal make raw: %s", err)
	}
	defer terminal.Restore(fd, state)

	w, h, err := terminal.GetSize(fd)
	if err != nil {
		return fmt.Errorf("terminal get size: %s", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	term := os.Getenv("TERM")
	if term == "" {
		term = "xterm-256color"
	}
	if err := session.RequestPty(term, h, w, modes); err != nil {
		return fmt.Errorf("session xterm: %s", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Shell(); err != nil {
		return fmt.Errorf("session shell: %s", err)
	}

	if err := session.Wait(); err != nil {
		if e, ok := err.(*ssh.ExitError); ok {
			switch e.ExitStatus() {
			case 130:
				return nil
			}
		}
		return fmt.Errorf("ssh: %s", err)
	}
	return nil
}
