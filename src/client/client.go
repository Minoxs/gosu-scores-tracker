package client

import (
	"log"
	"time"

	"osu-phantom/src/osu"
	"osu-phantom/src/utils"
)

var (
	CLIENT_ID     = utils.GetEnv("CLIENT_ID").Integer(0)
	CLIENT_SECRET = utils.GetEnv("CLIENT_SECRET").String()
)

type PhantomClient struct {
	Username string
	userID   int
	token    *osu.GuestToken
}

func Login(username string) (client *PhantomClient, err error) {
	client = &PhantomClient{Username: username}

	client.token, err = osu.GetGuestToken(CLIENT_ID, CLIENT_SECRET)
	if err != nil {
		return
	}
	client.userID, err = osu.GetUserID(client.token, client.Username)

	return
}

func (c *PhantomClient) Loop() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic: ", r)
		}
	}()

	var (
		lastSignal = (int64)(0)
		maxIdle    = (int64)((15 * time.Minute).Seconds())
		pool       = time.NewTimer(0 * time.Minute)
		ranking    = Ranking{}
	)
	defer pool.Stop()

	log.Println("Entering loop")
	defer func() {
		log.Println("Exiting loop")
	}()

	for {
		select {
		case <-pool.C:
			log.Println("Getting new scores...")
			var scores = c.GetRecentScores()
			log.Println("Score count: ", len(scores))

			// Check if there are new scores
			if len(scores) == 0 || scores[0].CreatedAt.Unix() == lastSignal {
				log.Println("No updates...")
				// Stop if idle for too long
				if time.Now().Unix()-lastSignal > maxIdle {
					return
				}
				// No new scores
				continue
			}

			// Add new scores to the ranks
			for _, score := range scores {
				if score.CreatedAt.Unix() <= lastSignal {
					break
				}

				log.Println("Potential new score: ", score.ID)
				ranking.AddScore(score)
			}

			// Update last signal
			lastSignal = scores[0].CreatedAt.Unix()

			// Reset timer
			pool.Reset(5 * time.Minute)
			log.Println(&ranking)
		}
	}
}

func (c *PhantomClient) GetRecentScores() osu.Scores {
	return osu.GetRecentScores(c.token, c.userID)
}

func (c *PhantomClient) GetBeatmapScores(beatmapID int) osu.Scores {
	return osu.GetBeatmapScores(c.token, c.userID, beatmapID)
}
