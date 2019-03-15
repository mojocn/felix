package flx

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sys/windows"
	"io"
	"log"
	"os"
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
	termWidth, termHeight := 80, 120
	//https://stackoverflow.com/questions/7949956/why-does-git-diff-on-windows-warn-that-the-terminal-is-not-fully-functional
	termType := "msys"
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fd := int(os.Stdin.Fd())
		state, err := terminal.MakeRaw(fd)
		if err != nil {
			return err
		}
		defer terminal.Restore(fd, state)
		fdOut := int(os.Stdout.Fd())
		termWidth, termHeight, err = terminal.GetSize(fdOut)
		if err != nil {
			return err
		}
	} else if isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		termType = "xterm"
	}

	err := t.Session.RequestPty(termType, termHeight, termWidth, ssh.TerminalModes{})
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
func makeCygwinRaw(fd int) {
	var raw = uint32(windows.ENABLE_PROCESSED_INPUT | windows.ENABLE_LINE_INPUT | windows.ENABLE_PROCESSED_OUTPUT)
	if err := windows.SetConsoleMode(windows.Handle(fd), raw); err != nil {
		log.Printf("cygwin set raw failed %s", err)
	}
}
