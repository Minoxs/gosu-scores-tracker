package main

import (
	"log"

	"osu-phantom/src/client"
)

func main() {
	var cli, err = client.Login("minoxs")

	if err != nil {
		log.Fatal(err)
	}

	cli.Loop()
}
