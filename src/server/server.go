package server

import (
	"errors"
	"fmt"
	"log"
	"net"

	"osu-phantom/src/config"
	"osu-phantom/src/osu"
	"osu-phantom/src/utils"
)

type PhantomServer struct {
	token    *osu.GuestToken
	listener net.Listener
}

func New() (srv *PhantomServer, err error) {
	srv = &PhantomServer{}
	srv.token, err = osu.GetGuestToken(config.ClientId, config.ClientSecret)
	return
}

func (s *PhantomServer) Start() error {
	// Create TCP socket
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", config.ServerPort))
	if err != nil {
		return err
	}

	// Enter server loop
	go s.Loop()
	return nil
}

func (s *PhantomServer) Stop() error {
	if s.listener == nil {
		return errors.New("phantom server already stopped")
	}
	return s.listener.Close()
}

func (s *PhantomServer) GetToken() *osu.GuestToken {
	return s.token
}

func (s *PhantomServer) Loop() {
	// Panic handler
	defer utils.PanicHandler("PhantomServer.Loop")
	// Close connection when loop ends
	defer func() {
		_ = s.listener.Close()
		s.listener = nil
	}()

	// Listen for connections
	for {
		var conn, err = s.listener.Accept()
		if err != nil {
			log.Println("PhantomServer.Loop : error accepting connection :", err)
			continue
		}
		go s.handleConnection(conn)
	}
}
