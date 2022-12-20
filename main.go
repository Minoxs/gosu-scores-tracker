package main

import (
	"log"
	"os"

	"github.com/kardianos/service"
	"osu-phantom/src"
	"osu-phantom/src/config"
)

func main() {
	var p = &phantom.Program{}
	var s, err = service.New(p, config.Service)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		arg := os.Args[1]

		if arg == "-install" || arg == "-i" {
			err = s.Install()
			if err != nil {
				log.Fatalf("Error installing service: %s", err)
			}
			log.Println("Installed Successfully")
			return
		}

		if arg == "-uninstall" || arg == "-u" {
			err = s.Uninstall()
			if err != nil {
				log.Fatalf("Error uninstalling service: %s", err)
			}
			log.Println("Uninstalled Successfully")
			return
		}
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
