package socks5ws

import "io"

type RelayTcp interface {
	io.Reader
	io.Writer
	io.Closer
}
