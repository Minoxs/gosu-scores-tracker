package client

import (
	"fmt"
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
	var maxIdle = (int64)((5 * time.Minute).Seconds())
	var lastSignal = time.Now().Unix()
	var pool = time.NewTimer(0 * time.Minute)
	defer pool.Stop()

	const RANK_SIZE = 100
	var ranking [RANK_SIZE]*osu.Score
	for i, _ := range ranking {
		ranking[i] = nil
	}

	fmt.Println("Entering loop")
	defer func() {
		fmt.Println("Exiting loop")
	}()

	for {
		select {
		case <-pool.C:
			fmt.Println("Getting new scores...")
			var scores = c.GetRecentScores()
			// Check for new scores
			if len(scores) == 0 || scores[0].CreatedAt.Unix() == lastSignal {
				fmt.Println("No updates...")
				// Stop if idle for too long
				if time.Now().Unix()-lastSignal > maxIdle {
					return
				}
				continue
			}
			lastSignal = scores[0].CreatedAt.Unix()

			for _, score := range scores {
				fmt.Println(score)

				// Add to ranking in the right place
				var pp = score.GetPP()
				for i, s := range ranking {
					// Ranking ended
					if s == nil {
						ranking[i] = &score
						break
					}
					// Check if score is better
					if pp > ranking[i].PP {
						// Shift all the rankings
						for j := RANK_SIZE - 2; j >= i; j-- {
							ranking[j+1] = ranking[j]
						}
						// Add score
						ranking[i] = &score
						break
					}
				}
			}

			// Print ranking
			var totalPP = 0.0
			for _, score := range ranking {
				if score == nil {
					break
				}

				fmt.Println("ID: ", score.ID, "PP: ", score.PP)
				totalPP += score.PP
			}
			fmt.Println("PP Total: ", totalPP)

			pool.Reset(1 * time.Minute)
		}
	}
}

func (c *PhantomClient) GetRecentScores() osu.Scores {
	return osu.GetRecentScores(c.token, c.userID)
}

func (c *PhantomClient) GetBeatmapScores(beatmapID int) osu.Scores {
	return osu.GetBeatmapScores(c.token, c.userID, beatmapID)
}
