package socks5ws

import (
	"log/slog"
	"net"
)

const (
	socks5Version             = 0x05
	socks5ReplyOkay           = 0x00
	socks5ReplyFail           = 0x01
	socks5ReplyReserved       = 0x00
	socks5CmdConnect          = 0x01
	socks5CmdBind             = 0x02
	socks5CmdUdpAssoc         = 0x03
	socks5AtypeIPv4           = 0x01
	socks5AtypeDomain         = 0x03
	socks5AtypeIPv6           = 0x04
	socks5UdpFragNotSupported = 0x00
	socks5UdpFragEnd          = 0x80

	bufferSize = 64 << 10
)

func socks5Response(conn net.Conn, ipv4 net.IP, port int, socks5OkayOrFail byte) {
	if socks5OkayOrFail != socks5ReplyOkay {
		ipv4 = net.IPv4zero
		port = 0
	}
	if ipv4 == nil {
		ipv4 = net.IPv4zero
	}
	if port < 0 || port > 65535 {
		port = 0
	}
	response := []byte{socks5Version, socks5OkayOrFail, socks5ReplyReserved, socks5AtypeIPv4, ipv4[0], ipv4[1], ipv4[2], ipv4[3], byte(port >> 8), byte(port & 0xff)}
	_, err := conn.Write(response)
	if err != nil {
		slog.Error("socks5 request rely failed to write", "err", err.Error())
	}
}
