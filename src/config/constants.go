package config

import "github.com/kardianos/service"

var (
	// Service configuration
	Service = &service.Config{
		Name:        "Phantom",
		DisplayName: "Phantom",
		Description: "Osu Phantom Server",
	}
)
