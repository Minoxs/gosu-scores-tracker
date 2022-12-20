package server

import (
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

func (s *PhantomServer) Loop() {

}

func (s *PhantomServer) Stop() error {
	return nil
}
