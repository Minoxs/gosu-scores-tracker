package server

import (
	"net"

	"osu-phantom/src/utils"
)

func (s *PhantomServer) handleConnection(conn net.Conn) {
	defer utils.PanicHandler("PhantomServer.handleConnection")
}
