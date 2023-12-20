package phantom

import (
	"fmt"
	"log"
	"osu-phantom/pkg/osu"
	"time"
)

type (
	AuthProvider interface {
		GetToken() *osu.GuestToken
	}

	Client struct {
		Username   string
		Provider   AuthProvider
		userID     int
		ranking    Ranking
		lastUpdate time.Time
	}
)

func Login(provider AuthProvider, username string) (client *Client, err error) {
	client = &Client{Username: username, Provider: provider}
	client.userID, err = osu.GetUserID(provider.GetToken(), client.Username)
	client.lastUpdate = time.Now()
	return
}

func (c *Client) Loop() {
	defer func() {
		if r := recover(); r != nil {
			c.log("Recovered from panic: ", r)
		}
	}()

	var (
		maxIdle = 15 * time.Minute
		pool    = time.NewTimer(0 * time.Minute)
	)
	defer pool.Stop()

	c.log("Entering loop")
	defer func() {
		c.log("Exiting loop")
	}()

	for {
		select {
		case <-pool.C:
			c.log("Getting new scores...")

			// Stop if idle for too long
			if !c.update() && time.Now().Sub(c.lastUpdate) > maxIdle {
				return
			}

			// Reset timer
			pool.Reset(5 * time.Minute)
			c.log(&c.ranking)
		}
	}
}

func (c *Client) update() bool {
	var scores = c.getRecentScores()
	c.log("Score count: ", len(scores))

	// Check if there are new scores
	if len(scores) == 0 || scores[0].CreatedAt.Equal(c.lastUpdate) {
		c.log("No updates...")
		return false
	}

	// Add new scores to the ranks
	for _, score := range scores {
		if score.CreatedAt.Before(c.lastUpdate) {
			break
		}

		c.log("Potential new score: ", score.ID)
		c.ranking.AddScore(score)
	}

	// Update last signal
	c.lastUpdate = scores[0].CreatedAt
	return true
}

func (c *Client) log(v ...any) {
	var l = "PhantomClient." + c.Username + " : " + fmt.Sprint(v...)
	log.Println(l)
}

func (c *Client) getRecentScores() osu.Scores {
	return osu.GetRecentScores(c.Provider.GetToken(), c.userID)
}

func (c *Client) getBeatmapScores(beatmapID int) osu.Scores {
	return osu.GetBeatmapScores(c.Provider.GetToken(), c.userID, beatmapID)
}
