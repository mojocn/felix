package socks5ws

type VlessCmd byte
type VlessAddrType byte

const (
	vlessCmdTcp      VlessCmd      = 0x01
	vlessCmdUdp      VlessCmd      = 0x02
	vlessCmdMux      VlessCmd      = 0x03
	vlessAtypeIPv4   VlessAddrType = 0x01
	vlessAtypeDomain VlessAddrType = 0x02
	vlessAtypeIPv6   VlessAddrType = 0x03

	socks5Version       = 0x05
	socks5ReplySuccess  = 0x00
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
	socks5ReplyBytesSuccess = []byte{socks5Version, socks5ReplySuccess, socks5ReplyReserved, socks5AtypeIPv4, 0, 0, 0, 0, 0, 0}
	socks5ReplyBytesFailed  = []byte{socks5Version, socks5ReplyFail, socks5ReplyReserved, socks5AtypeIPv4, 0, 0, 0, 0, 0, 0}
)
