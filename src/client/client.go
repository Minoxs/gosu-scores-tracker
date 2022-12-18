package client

import (
	"osu-phantom/src/osu"
	"osu-phantom/src/utils"
)

var (
	CLIENT_ID     = utils.GetEnv("client_id").Integer(0)
	CLIENT_SECRET = utils.GetEnv("client_secret").String()
)

type PhantomClient struct {
	Username string
	userID   int
	token    *osu.GuestToken
}

func (c *PhantomClient) Login() {
	var err error

	c.token = osu.GetGuestToken(CLIENT_ID, CLIENT_SECRET)
	if c.token == nil {
		panic("error getting token")
	}
	c.userID, err = osu.GetUserID(c.token, c.Username)
	if err != nil {
		panic(err)
	}
}

func (c *PhantomClient) GetRecentScores() osu.Scores {
	return osu.GetRecentScores(c.token, c.userID)
}

func (c *PhantomClient) GetBeatmapScores(beatmapID int) osu.Scores {
	return osu.GetBeatmapScores(c.token, c.userID, beatmapID)
}
