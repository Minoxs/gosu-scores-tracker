package server

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/minoxs/gommunication"
	"osu-phantom/src/utils"
)

func parseMessage(buf io.Reader, h *gommunication.Header) (proc processable, err error) {
	switch h.ID {
	}
	return
}

func (s *PhantomServer) checkLogin(conn net.Conn) (err error) {
	var login = gommunication.Message[ServerLogin]{}
	err = login.FromStream(conn)
	if err != nil {
		return
	}
	log.Println("Login received: ", login.Body.Name)
	return
}

func (s *PhantomServer) handleConnection(conn net.Conn) {
	defer utils.PanicHandler("PhantomServer.handleConnection")

	var (
		err  error
		hdr  gommunication.Header
		proc processable
	)

	defer func() {
		log.Println("Closed connection : " + err.Error())
	}()

	// Check for login first
	err = s.checkLogin(conn)
	if err != nil {
		return
	}
	_ = conn.SetDeadline(time.Now().Add(5 * time.Minute))

	for {
		// Receive header
		err = hdr.FromStream(conn)
		if err != nil {
			return
		}

		// Extend connection life
		_ = conn.SetDeadline(time.Now().Add(5 * time.Minute))

		// Parse it
		proc, err = parseMessage(conn, &hdr)
		if err == gommunication.MissingEOM {
			log.Println("Attempting to flush message")
			err = gommunication.FlushMessage(conn)
		}
		if err != nil {
			return
		}
		if proc == nil {
			continue
		}

		// Process request
		proc.Process(conn)
	}
}
