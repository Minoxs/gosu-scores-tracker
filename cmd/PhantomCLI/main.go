package main

import (
	"log"
	"os"
	"time"

	"github.com/minoxs/osu-phantom/pkg/osu"
	"github.com/minoxs/osu-phantom/pkg/phantom"
)

func FatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	var (
		credentials osu.Credentials
		username    string
	)

	// Parse Arguments
	ParseArgs(os.Args).
		Int("ClientID", "id", &credentials.ClientID).Required().
		String("ClientSecret", "s", &credentials.ClientSecret).Required().
		String("Username", "u", &username).Required()

	var (
		err    error
		token  *osu.GuestToken
		client *phantom.Client
	)

	token, err = osu.NewClient(0).GetGuestToken(credentials)
	FatalError(err)

	var provider = &phantom.DefaultProvider{Token: token}

	client, err = phantom.Login(provider, username, time.Now())
	FatalError(err)

	client.KeepUpdated(1*time.Minute, 30*time.Minute)
}
