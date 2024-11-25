package shadowos

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net"
)

type Socks5Request struct {
	id          string
	socks5Cmd   byte
	socks5Atyp  byte
	dstAddr     []byte
	dstPort     []byte
	CountryCode string //iso country code
}

func parseSocks5Request(data []byte, geo *GeoIP) (*Socks5Request, error) {
	id := uuid.NewString()
	info := &Socks5Request{id: id}

	if data[0] != socks5Version {
		return nil, fmt.Errorf("unsupported SOCKS version: %d", data[0])
	}
	if data[1] == socks5CmdConnect {
		info.socks5Cmd = socks5CmdConnect
	} else if data[1] == socks5CmdUdpAssoc {
		info.socks5Cmd = socks5CmdUdpAssoc
	} else {
		//BIND is not supported
		return nil, fmt.Errorf("unsupported command: %d", data[1])
	}
	if data[2] != socks5ReplyReserved {
		return nil, fmt.Errorf("RSV must be 0x00")
	}
	if data[3] == socks5AtypeIPv4 {
		if len(data) < 10 {
			return nil, fmt.Errorf("request too short for atyp IPv4")
		}
		info.socks5Atyp = socks5AtypeIPv4
		info.dstAddr = data[4:8]
		info.dstPort = data[8:10]
	} else if data[3] == socks5AtypeDomain {
		if len(data) < 5 {
			return nil, fmt.Errorf("request too short for atyp Domain")
		}
		addrLen := int(data[4])
		info.socks5Atyp = socks5AtypeDomain
		info.dstAddr = data[5 : 5+addrLen]
		info.dstPort = data[5+addrLen : 5+addrLen+2]
	} else if data[3] == socks5AtypeIPv6 {
		if len(data) < 22 {
			return nil, fmt.Errorf("request too short for atyp IPv6")
		}
		info.socks5Atyp = socks5AtypeIPv6
		info.dstAddr = data[4:20]
		info.dstPort = data[20:22]
	} else {
		return nil, fmt.Errorf("unsupported address type: %d", data[3])
	}
	code, err := geo.country(info.host())
	if err != nil {
		info.Logger().Error("failed to get country code", "err", err.Error())
	} else {
		info.CountryCode = code
	}
	return info, nil
}

func (s Socks5Request) host() string {
	addr := ""
	if s.socks5Atyp == socks5AtypeIPv4 || s.socks5Atyp == socks5AtypeIPv6 {
		addr = net.IP(s.dstAddr).String()
	} else if s.socks5Atyp == socks5AtypeDomain {
		addr = string(s.dstAddr)
	} else {
		addr = string(s.dstAddr)
	}
	return addr
}
func (s Socks5Request) addr() string {
	return fmt.Sprintf("%s:%s", s.host(), s.port())
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

func (s Socks5Request) Logger() *slog.Logger {
	return slog.With("reqId", s.id, "cmd", s.cmd(), "atyp", s.aType(), "host", s.host(), "port", s.port(), "country", s.CountryCode)
}
func (s Socks5Request) String() string {
	return fmt.Sprintf("socks5Cmd: %v, socks5Atyp: %v, dstAddr: %v, dstPort: %v, country: %s", s.cmd(), s.aType(), s.host(), s.port())
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
