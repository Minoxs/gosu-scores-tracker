package main

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

	program struct {
		server thread
	}
)

func (p *program) Start(_ service.Service) (err error) {
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

func (p *program) Stop(_ service.Service) error {
	return p.server.Stop()
}
