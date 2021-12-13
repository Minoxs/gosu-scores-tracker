package main

import (
	"fmt"

	"osu-phantom/client"
)

func main() {
	Client := &client.PhantomClient{
		Username: "minoxs",
	}
	Client.Login()
	fmt.Println("Recent Scores: ", Client.GetRecentScores())
}
