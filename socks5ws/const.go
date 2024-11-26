package socks5ws

import "net"

const (
	socks5Version       = 0x05
	socks5ReplyOkay     = 0x00
	socks5ReplyFail     = 0x01
	socks5ReplyReserved = 0x00
	socks5CmdConnect    = 0x01
	socks5CmdBind       = 0x02
	socks5CmdUdpAssoc   = 0x03
	socks5AtypeIPv4     = 0x01
	socks5AtypeDomain   = 0x03
	socks5AtypeIPv6     = 0x04
)

var (
	socks5ReplyBytesOkay = []byte{socks5Version, socks5ReplyOkay, socks5ReplyReserved, socks5AtypeIPv4, net.IPv4zero[0], net.IPv4zero[1], net.IPv4zero[2], net.IPv4zero[3], 0, 0}
	socks5ReplyBytesFail = []byte{socks5Version, socks5ReplyFail, socks5ReplyReserved, socks5AtypeIPv4, net.IPv4zero[0], net.IPv4zero[1], net.IPv4zero[2], net.IPv4zero[3], 0, 0}
)
