package shadowos

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Socks5Request struct {
	id         string
	socks5Cmd  byte
	socks5Atyp byte
	dstAddr    []byte
	dstPort    []byte
}

func NewSocks5Request(id string) *Socks5Request {
	if id == "" {
		id = uuid.NewString()
	}
	return &Socks5Request{id: id}
}

func (s Socks5Request) addr() string {
	addr := ""
	if s.socks5Atyp == socks5AtypeIPv4 {
		addr = fmt.Sprintf("%d.%d.%d.%d", s.dstAddr[0], s.dstAddr[1], s.dstAddr[2], s.dstAddr[3])
	} else if s.socks5Atyp == socks5AtypeIPv6 {
		addr = fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x", s.dstAddr[0], s.dstAddr[1], s.dstAddr[2], s.dstAddr[3], s.dstAddr[4], s.dstAddr[5], s.dstAddr[6], s.dstAddr[7])
	} else if s.socks5Atyp == socks5AtypeDomain {
		addr = string(s.dstAddr)
	} else {
		addr = string(s.dstAddr)
	}
	return addr
}
func (s Socks5Request) cmd() string {
	cmd := "unknown"
	if s.socks5Cmd == socks5CmdConnect {
		cmd = "connect"
	} else if s.socks5Cmd == socks5CmdUdpAssoc {
		cmd = "udp"
	} else if s.socks5Cmd == socks5CmdBind {
		cmd = "bind"
	}
	return cmd
}
func (s Socks5Request) aType() string {
	return fmt.Sprintf("%v", s.socks5Atyp)
}

func (s Socks5Request) port() string {
	port := int(s.dstPort[0])<<8 + int(s.dstPort[1])
	return fmt.Sprintf("%v", port)
}

func (s Socks5Request) Logger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"reqId": s.id,
		"cmd":   s.cmd(),
		"atyp":  s.aType(),
		"addr":  s.addr(),
		"port":  s.port(),
	})
}
func (s Socks5Request) String() string {
	return fmt.Sprintf("socks5Cmd: %v, socks5Atyp: %v, dstAddr: %v, dstPort: %v", s.cmd(), s.aType(), s.addr(), s.addr())
}

func (s Socks5Request) addressBytes() []byte {
	if s.socks5Atyp == socks5AtypeDomain {
		return append([]byte{byte(len(s.dstAddr))}, s.dstAddr...)
	}
	return s.dstAddr
}

func (s Socks5Request) vlessHeaderTcp(uuid [16]byte) ([]byte, error) {
	addrBytes := s.addressBytes()
	//https://xtls.github.io/development/protocols/vless.html
	headerBytes := make([]byte, 0, 1+16+1+1+2+1+len(addrBytes))

	headerBytes = append(headerBytes, 0x01)       // version
	headerBytes = append(headerBytes, uuid[:]...) //16 bytes of UUID
	headerBytes = append(headerBytes, 0x00)       // additional info length M

	//1 byte of command
	if s.socks5Cmd == socks5CmdUdpAssoc {
		headerBytes = append(headerBytes, byte(vlessCmdUdp))
	} else if s.socks5Cmd == socks5CmdConnect {
		headerBytes = append(headerBytes, byte(vlessCmdTcp))
	} else {
		return nil, fmt.Errorf("unsupported command: %d", s.socks5Cmd)
	}

	headerBytes = append(headerBytes, s.dstPort...) //2 bytes of port

	//1 byte of address type
	if s.socks5Atyp == socks5AtypeIPv4 {
		headerBytes = append(headerBytes, byte(vlessAtypeIPv4))
	} else if s.socks5Atyp == socks5AtypeIPv6 {
		headerBytes = append(headerBytes, byte(vlessAtypeIPv6))
	} else if s.socks5Atyp == socks5AtypeDomain {
		headerBytes = append(headerBytes, byte(vlessAtypeDomain))
	} else {
		return nil, fmt.Errorf("unsupported address type: %d", s.socks5Atyp)
	}

	headerBytes = append(headerBytes, addrBytes...) //n bytes of address
	return headerBytes, nil
}

func (s Socks5Request) vlessHeaderUdp(uuid [16]byte) ([]byte, error) {
	addrBytes := s.addressBytes()
	//https://xtls.github.io/development/protocols/vless.html
	headerBytes := make([]byte, 0, 1+16+1+1+2+1+len(addrBytes))

	headerBytes = append(headerBytes, 0x01)       // version
	headerBytes = append(headerBytes, uuid[:]...) //16 bytes of UUID
	headerBytes = append(headerBytes, 0x00)       // additional info length M

	//1 byte of command
	if s.socks5Cmd == socks5CmdUdpAssoc {
		headerBytes = append(headerBytes, byte(vlessCmdUdp))
	} else if s.socks5Cmd == socks5CmdConnect {
		headerBytes = append(headerBytes, byte(vlessCmdTcp))
	} else {
		return nil, fmt.Errorf("unsupported command: %d", s.socks5Cmd)
	}

	headerBytes = append(headerBytes, s.dstPort...) //2 bytes of port

	//1 byte of address type
	if s.socks5Atyp == socks5AtypeIPv4 {
		headerBytes = append(headerBytes, byte(vlessAtypeIPv4))
	} else if s.socks5Atyp == socks5AtypeIPv6 {
		headerBytes = append(headerBytes, byte(vlessAtypeIPv6))
	} else if s.socks5Atyp == socks5AtypeDomain {
		headerBytes = append(headerBytes, byte(vlessAtypeDomain))
	} else {
		return nil, fmt.Errorf("unsupported address type: %d", s.socks5Atyp)
	}

	headerBytes = append(headerBytes, addrBytes...) //n bytes of address
	return headerBytes, nil
}
