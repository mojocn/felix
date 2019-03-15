package flx

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/dejavuzhou/felix/model"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

const sudoPrefix, sudoSuffix = "[sudo] password for ", ": "
const sudoPrefixLen = len(sudoPrefix)

type SSHTerminal struct {
	Session            *ssh.Session
	exitMsg            string
	stdout             io.Reader
	stdin              io.Writer
	stderr             io.Reader
	Password           string
	LoginUser          string
	EnableSudoPassword bool
}

func RunSshTerminal(h *model.Machine, sudoMode bool) error {
	client, err := NewSshClient(h)
	if err != nil {
		return err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	if h.Password == "" {
		sudoMode = false
	}
	s := SSHTerminal{
		Session:            session,
		Password:           h.Password,
		LoginUser:          h.User,
		EnableSudoPassword: sudoMode,
	}
	return s.interactiveSession()
}

func enableSudoPassword(t *SSHTerminal) {
	var (
		line string
		r    = bufio.NewReader(t.stdout)
	)
	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		line += string(b)
		os.Stdout.Write([]byte{b})
		if b == byte('\n') {
			line = ""
			continue
		}
		if len(line) >= sudoPrefixLen && strings.HasPrefix(line, sudoPrefix) && strings.HasSuffix(line, sudoSuffix) && strings.Contains(line, t.LoginUser) {
			_, err = t.stdin.Write([]byte(t.Password + "\n"))
			if err != nil {
				break
			}
			color.Green("\r\nFelix has automatically input password for %s", color.BlueString(t.LoginUser))
		}
	}
}
