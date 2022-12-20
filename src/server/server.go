package server

import "errors"

type PhantomServer struct {
}

func New() (srv *PhantomServer, err error) {
	return nil, errors.New("IMPLEMENT")
}

func (s *PhantomServer) Start() {

}

func (s *PhantomServer) Stop() error {
	return nil
}
