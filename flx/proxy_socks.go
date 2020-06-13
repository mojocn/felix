package flx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/libragen/felix/model"
)

type socksDialer interface {
	DialTCP(net string, laddr, raddr *net.TCPAddr) (net.Conn, error)
}

func RunSocksProxy(h *model.Machine, port int) error {
	conn, err := NewSshClient(h)
	if err != nil {
		return err
	}
	defer conn.Close()

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("unable to listen on SOCKS port [%s]: %v", addr, err)
	}
	defer l.Close()
	fmt.Printf("socks porxy is up on 127.0.0.1:%d", port)
	for {
		c, err := l.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept incoming SOCKS connection: %v", err)
		}
		go handleConn(c.(*net.TCPConn), conn)
	}
	connections.Wait()
	return nil
}

var connections = new(sync.WaitGroup)

func handleConn(local *net.TCPConn, dialer socksDialer) {
	connections.Add(1)
	defer local.Close()
	defer connections.Done()

	// SOCKS does not include a length in the header, so take
	// a punt that each request will be readable in one go.
	buf := make([]byte, 256)
	n, err := local.Read(buf)
	if err != nil || n < 2 {
		log.Printf("[%s] unable to read SOCKS header: %v", local.RemoteAddr(), err)
		return
	}
	buf = buf[:n]

	switch version := buf[0]; version {
	case 4:
		switch command := buf[1]; command {
		case 1:
			port := binary.BigEndian.Uint16(buf[2:4])
			ip := net.IP(buf[4:8])
			addr := &net.TCPAddr{IP: ip, Port: int(port)}
			buf := buf[8:]
			i := bytes.Index(buf, []byte{0})
			if i < 0 {
				log.Printf("[%s] unable to locate SOCKS4 user", local.RemoteAddr())
				return
			}
			user := buf[:i]
			log.Printf("[%s] incoming SOCKS4 TCP/IP stream connection, user=%q, raddr=%s", local.RemoteAddr(), user, addr)
			remote, err := dialer.DialTCP("tcp4", local.RemoteAddr().(*net.TCPAddr), addr)
			if err != nil {
				log.Printf("[%s] unable to connect to remote host: %v", local.RemoteAddr(), err)
				local.Write([]byte{0, 0x5b, 0, 0, 0, 0, 0, 0})
				return
			}
			local.Write([]byte{0, 0x5a, 0, 0, 0, 0, 0, 0})
			transfer(local, remote)
		default:
			log.Printf("[%s] unsupported command, closing connection", local.RemoteAddr())
		}
	case 5:
		authlen, buf := buf[1], buf[2:]
		auths, buf := buf[:authlen], buf[authlen:]
		if !bytes.Contains(auths, []byte{0}) {
			log.Printf("[%s] unsuported SOCKS5 authentication method", local.RemoteAddr())
			local.Write([]byte{0x05, 0xff})
			return
		}
		local.Write([]byte{0x05, 0x00})
		buf = make([]byte, 256)
		n, err := local.Read(buf)
		if err != nil {
			log.Printf("[%s] unable to read SOCKS header: %v", local.RemoteAddr(), err)
			return
		}
		buf = buf[:n]
		switch version := buf[0]; version {
		case 5:
			switch command := buf[1]; command {
			case 1:
				buf = buf[3:]
				switch addrtype := buf[0]; addrtype {
				case 1:
					if len(buf) < 8 {
						log.Printf("[%s] corrupt SOCKS5 TCP/IP stream connection request", local.RemoteAddr())
						local.Write([]byte{0x05, 0x07, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
						return
					}
					ip := net.IP(buf[1:5])
					port := binary.BigEndian.Uint16(buf[5:6])
					addr := &net.TCPAddr{IP: ip, Port: int(port)}
					log.Printf("[%s] incoming SOCKS5 TCP/IP stream connection, raddr=%s", local.RemoteAddr(), addr)
					remote, err := dialer.DialTCP("tcp", local.RemoteAddr().(*net.TCPAddr), addr)
					if err != nil {
						log.Printf("[%s] unable to connect to remote host: %v", local.RemoteAddr(), err)
						local.Write([]byte{0x05, 0x04, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
						return
					}
					local.Write([]byte{0x05, 0x00, 0x00, 0x01, ip[0], ip[1], ip[2], ip[3], byte(port >> 8), byte(port)})
					transfer(local, remote)
				case 3:
					addrlen, buf := buf[1], buf[2:]
					name, buf := buf[:addrlen], buf[addrlen:]
					ip, err := net.ResolveIPAddr("ip", string(name))
					if err != nil {
						log.Printf("[%s] unable to resolve IP address: %q, %v", local.RemoteAddr(), name, err)
						local.Write([]byte{0x05, 0x04, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
						return
					}
					port := binary.BigEndian.Uint16(buf[:2])
					addr := &net.TCPAddr{IP: ip.IP, Port: int(port)}
					remote, err := dialer.DialTCP("tcp", local.RemoteAddr().(*net.TCPAddr), addr)
					if err != nil {
						log.Printf("[%s] unable to connect to remote host: %v", local.RemoteAddr(), err)
						local.Write([]byte{0x05, 0x04, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
						return
					}
					local.Write([]byte{0x05, 0x00, 0x00, 0x01, addr.IP[0], addr.IP[1], addr.IP[2], addr.IP[3], byte(port >> 8), byte(port)})
					transfer(local, remote)

				default:
					log.Printf("[%s] unsupported SOCKS5 address type: %d", local.RemoteAddr(), addrtype)
					local.Write([]byte{0x05, 0x08, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
				}
			default:
				log.Printf("[%s] unknown SOCKS5 command: %d", local.RemoteAddr(), command)
				local.Write([]byte{0x05, 0x07, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
			}
		default:
			log.Printf("[%s] unnknown version after SOCKS5 handshake: %d", local.RemoteAddr(), version)
			local.Write([]byte{0x05, 0x07, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
		}
	default:
		log.Printf("[%s] unknown SOCKS version: %d", local.RemoteAddr(), version)
	}
}

func transfer(in, out net.Conn) {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	f := func(in, out net.Conn, wg *sync.WaitGroup) {
		n, err := io.Copy(out, in)
		log.Printf("xfer done: in=%v\tout=%v\ttransfered=%d\terr=%v", in.RemoteAddr(), out.RemoteAddr(), n, err)
		if conn, ok := in.(*net.TCPConn); ok {
			conn.CloseWrite()
		}
		if conn, ok := out.(*net.TCPConn); ok {
			conn.CloseRead()
		}
		wg.Done()
	}
	go f(in, out, wg)
	f(out, in, wg)
	wg.Wait()
	out.Close()
}
