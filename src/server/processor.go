package server

import "net"

type (
	processable interface {
		Process(net.Conn)
	}
)
