package main

import (
	"osu-phantom/src/server"
)

func (p *program) initRequired() (err error) {
	p.server, err = server.New()
	if err != nil {
		return err
	}

	return
}

func (p *program) initModules() {
	go p.server.Loop()
}
