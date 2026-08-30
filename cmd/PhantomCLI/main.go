package main

import (
	"log"
	"os"
	"time"

	"github.com/minoxs/gosu-api/pkg/gosu"
	"github.com/minoxs/gosu-scores-tracker/pkg/tracker"
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
		credentials gosu.Credentials
		username    string
	)

	// Parse Arguments
	ParseArgs(os.Args).
		Int("ClientID", "id", &credentials.ClientID).Required().
		String("ClientSecret", "s", &credentials.ClientSecret).Required().
		String("Username", "u", &username).Required()

	var (
		err    error
		token  *gosu.GuestToken
		client *tracker.Client
	)

	token, err = gosu.NewClient(0).GetGuestToken(credentials)
	FatalError(err)

	var provider = &tracker.DefaultProvider{Token: token}

	client, err = tracker.Login(provider, username, time.Now())
	FatalError(err)

	client.KeepUpdated(1*time.Minute, 30*time.Minute)
}
