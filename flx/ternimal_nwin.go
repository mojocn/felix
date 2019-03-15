// +build !windows

package flx

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"runtime"
	"time"
)

func (t *SSHTerminal) interactiveSession() error {

	defer func() {
		if t.exitMsg == "" {
			fmt.Fprintln(os.Stdout, "[Felix]: the connection was closed on the remote side on ", time.Now().Format(time.RFC822))
		} else {
			fmt.Fprintln(os.Stdout, t.exitMsg)
		}
	}()
	fd := int(os.Stdin.Fd())
	if !terminal.IsTerminal(fd) {
		osName := runtime.GOOS
		return fmt.Errorf("%s fd %d is not a terminal,can't create pty of ssh", osName, fd)
	}
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, state)

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}
	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}

	err = t.Session.RequestPty(termType, termHeight, termWidth, ssh.TerminalModes{})
	if err != nil {
		return err
	}

	t.updateTerminalSize()

	t.stdin, err = t.Session.StdinPipe()
	if err != nil {
		return err
	}
	t.stdout, err = t.Session.StdoutPipe()
	if err != nil {
		return err
	}
	t.stderr, err = t.Session.StderrPipe()

	go io.Copy(os.Stderr, t.stderr)
	if t.EnableSudoPassword {
		go enableSudoPassword(t)
	} else {
		go io.Copy(os.Stdout, t.stdout)
	}
	go io.Copy(t.stdin, os.Stdin)

	err = t.Session.Shell()
	if err != nil {
		return err
	}
	return t.Session.Wait()
}
