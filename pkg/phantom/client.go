package client

import (
	"fmt"
	"log"
	"time"

	"osu-phantom/src/osu"
)

type (
	AuthProvider interface {
		GetToken() *osu.GuestToken
	}

	PhantomClient struct {
		Username   string
		Provider   AuthProvider
		userID     int
		ranking    Ranking
		lastUpdate int64
	}
)

func Login(provider AuthProvider, username string) (client *PhantomClient, err error) {
	client = &PhantomClient{Username: username, Provider: provider}
	client.userID, err = osu.GetUserID(provider.GetToken(), client.Username)
	return
}

func (c *PhantomClient) Loop() {
	defer func() {
		if r := recover(); r != nil {
			c.log("Recovered from panic: ", r)
		}
	}()

	var (
		maxIdle = (int64)((15 * time.Minute).Seconds())
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
			if !c.update() && time.Now().Unix()-c.lastUpdate > maxIdle {
				return
			}

			// Reset timer
			pool.Reset(5 * time.Minute)
			c.log(&c.ranking)
		}
	}
}

func (c *PhantomClient) update() bool {
	var scores = c.getRecentScores()
	c.log("Score count: ", len(scores))

	// Check if there are new scores
	if len(scores) == 0 || scores[0].CreatedAt.Unix() == c.lastUpdate {
		c.log("No updates...")
		return false
	}

	// Add new scores to the ranks
	for _, score := range scores {
		if score.CreatedAt.Unix() <= c.lastUpdate {
			break
		}

		c.log("Potential new score: ", score.ID)
		c.ranking.AddScore(score)
	}

	// Update last signal
	c.lastUpdate = scores[0].CreatedAt.Unix()
	return true
}

func (c *PhantomClient) log(v ...any) {
	var l = "PhantomClient." + c.Username + " : " + fmt.Sprint(v...)
	log.Println(l)
}

func (c *PhantomClient) getRecentScores() osu.Scores {
	return osu.GetRecentScores(c.Provider.GetToken(), c.userID)
}

func (c *PhantomClient) getBeatmapScores(beatmapID int) osu.Scores {
	return osu.GetBeatmapScores(c.Provider.GetToken(), c.userID, beatmapID)
}
