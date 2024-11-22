package shadowos

import "fmt"

type Socks5Request struct {
	socks5Cmd  byte
	socks5Atyp byte
	dstAddr    []byte
	dstPort    []byte
}

func (s Socks5Request) String() string {
	return fmt.Sprintf("socks5Cmd: %v, socks5Atyp: %v, dstAddr: %v, dstPort: %v", s.socks5Cmd, s.socks5Atyp, s.dstAddr, s.dstPort)
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
