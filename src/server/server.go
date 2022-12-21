package server

import (
	"log"

	"osu-phantom/src/client"
	"osu-phantom/src/config"
	"osu-phantom/src/osu"
)

type PhantomServer struct {
	token *osu.GuestToken
}

func New() (srv *PhantomServer, err error) {
	srv = &PhantomServer{}
	srv.token, err = osu.GetGuestToken(config.CLIENT_ID, config.CLIENT_SECRET)
	return
}

func (p *PhantomServer) GetToken() *osu.GuestToken {
	return p.token
}

func (s *PhantomServer) Loop() {
	var test, err = client.Login(s, "minoxs")
	if err != nil {
		log.Println(err)
	}
	go test.Loop()
}

func (s *PhantomServer) Stop() error {
	return nil
}
