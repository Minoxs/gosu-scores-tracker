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

	scores := Client.GetRecentScores()
	fmt.Println("Recent Scores: ", scores)
}
