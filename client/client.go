package client

import (
	"osu-phantom/osu"
)

const CLIENT_ID = 0
const CLIENT_SECRET = ""

type PhantomClient struct {
	Username string
	userID   int
	token    *osu.GuestToken
}

func (c *PhantomClient) Login() {
	var err error

	c.token = osu.GetGuestToken(CLIENT_ID, CLIENT_SECRET)
	c.userID, err = osu.GetUserID(c.token, c.Username)
	if err != nil {
		panic(err)
	}
}

func (c *PhantomClient) GetRecentScores() []osu.Score {
	return osu.GetRecentScores(c.token, c.userID)
}
