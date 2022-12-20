package phantom

import (
	"errors"
	"log"

	"github.com/kardianos/service"
)

type (
	module interface {
		Stop() error
	}

	thread interface {
		module
		Loop()
	}

	Program struct {
		server thread
	}
)

func (p *Program) Start(_ service.Service) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			err = errors.New("panic while initializing")
		}
	}()

	err = p.initRequired()
	if err != nil {
		return err
	}
	go p.initModules()

	return nil
}

func (p *Program) Stop(_ service.Service) error {
	return p.server.Stop()
}
