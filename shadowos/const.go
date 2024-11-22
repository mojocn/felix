package shadowos

type VlessCmd byte
type VlessAddrType byte

const (
	socksVersion       = 0x05
	cmdUDPAssociate    = 0x03
	socks5ReplySuccess = 0x00
	socks5ReplyFail    = 0x01

	vlessCmdTcp VlessCmd = 0x01
	vlessCmdUdp VlessCmd = 0x02
	vlessCmdMux VlessCmd = 0x03

	vlessAtypeIPv4   VlessAddrType = 0x01
	vlessAtypeDomain VlessAddrType = 0x02
	vlessAtypeIPv6   VlessAddrType = 0x03

	SOCKS5VERSION     = 0x05
	socks5CmdConnect  = 0x01
	socks5CmdBind     = 0x02
	socks5CmdUdpAssoc = 0x03
	socks5AtypeIPv4   = 0x01
	socks5AtypeDomain = 0x03
	socks5AtypeIPv6   = 0x04
)
