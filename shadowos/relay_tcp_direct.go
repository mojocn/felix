package shadowos

import "net"

var _ RelayTcp = (*RelayTcpDirect)(nil)

type RelayTcpDirect struct {
	conn net.Conn
}

func NewRelayTcpDirect(req *Socks5Request) (*RelayTcpDirect, error) {
	conn, err := net.Dial("tcp", req.addr())
	if err != nil {
		return nil, err
	}
	return &RelayTcpDirect{conn: conn}, nil
}

func (r *RelayTcpDirect) Read(p []byte) (n int, err error) {
	return r.conn.Read(p)
}

func (r *RelayTcpDirect) Write(p []byte) (n int, err error) {
	return r.conn.Write(p)
}

func (r *RelayTcpDirect) Close() error {
	return r.conn.Close()
}
