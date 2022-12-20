package main

import (
	"errors"
	"log"

	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(_ service.Service) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			err = errors.New("panic while initializing")
		}
	}()

	err = initRequired()
	if err != nil {
		return err
	}
	go initModules()

	return nil
}

func (p *program) Stop(_ service.Service) error {
	return nil
}
