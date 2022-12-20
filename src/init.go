package phantom

import (
	"osu-phantom/src/server"
)

func (p *Program) initRequired() (err error) {
	p.server, err = server.New()
	if err != nil {
		return err
	}

	return
}

func (p *Program) initModules() {
	go p.server.Loop()
}
