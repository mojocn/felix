package flx

import (
	"fmt"
	"io"
	"net"

	"github.com/mojocn/felix/model"
	"golang.org/x/crypto/ssh"
)

// sshProxy is a struct for SSH tunnel configuration
type sshProxy struct {
	localAddr  string
	remoteAddr string
	client     *ssh.Client
}

func RunProxy(h *model.Machine, localAddr, remoteAddr string) error {
	client, err := NewSshClient(h)
	if err != nil {
		return err
	}
	defer client.Close()
	tunnel := sshProxy{localAddr, remoteAddr, client}
	return tunnel.start()
}

// Start Method to start a local server and forward connection to the remote one.
func (tunnel *sshProxy) start() error {
	listener, err := net.Listen("tcp", tunnel.localAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}
func (tunnel *sshProxy) forward(localConn net.Conn) {
	remoteConn, err := tunnel.client.Dial("tcp", tunnel.remoteAddr)
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}
	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}
