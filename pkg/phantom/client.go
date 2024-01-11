package phantom

import (
	"errors"
	"fmt"
	"github.com/minoxs/osu-phantom/pkg/osu"
	"github.com/minoxs/osu-phantom/pkg/osu/player"
	"log"
	"sync"
	"time"
)

// TODO CONFIGURE START DATE
// TODO AVOID RECALCULATING PP FOR SCORE (CHECK SCORE ID)
// TODO USE SLOG IN OSU PHANTOM

type (
	// AuthProvider is required for requests which require OAuth authorization
	AuthProvider interface {
		GetToken() *osu.GuestToken
	}

	// Client handles the tracking of a user's scores
	Client struct {
		UserID   int
		Username string
		Provider AuthProvider

		ranking    player.Ranking
		LastUpdate time.Time

		lock sync.Mutex
	}
)

var ErrUserNotFound = errors.New("user not found")

// Login will look for the user and return a phantom client with the given user.
// Returns ErrUserNotFound if user is not found, or some HTTP error if failed to fetch user.
// Call Client.KeepUpdated to keep rankings constantly updated, or manually update with Client.Update.
func Login(provider AuthProvider, username string) (client *Client, err error) {
	client = &Client{Username: username, Provider: provider}
	client.UserID, err = osu.GetUserID(provider.GetToken(), client.Username)
	client.LastUpdate = time.Now()

	// API only returns error if request failed.
	// If user was not found, return ErrUserNotFound.
	if err == nil && client.UserID == 0 {
		err = ErrUserNotFound
	}

	// Unset client pointer if function errored.
	// This is only done to make the function fit go standards.
	if err != nil {
		client = nil
	}

	return
}

// KeepUpdated will fetch new scores from the API in the interval configured.
// Will stop routine after maxIdle without new scores.
// Returns a function that will prematurely stop the routine if called.
func (c *Client) KeepUpdated(checkInterval time.Duration, maxIdle time.Duration) {
	c.log("Running KeepUpdated")
	defer func() {
		if r := recover(); r != nil {
			c.log("Recovered from panic: ", r)
		}
		c.log("Stopping KeepUpdated routine")
	}()

	var interval = time.NewTimer(0)
	defer interval.Stop()

	for {
		select {
		case <-interval.C:
			c.log("Getting new scores...")

			// Stop if idle for too long
			if !c.Update() && time.Now().Sub(c.LastUpdate) > maxIdle {
				return
			}

			// Reset timer
			interval.Reset(checkInterval)
			c.log(&c.ranking)
		}
	}
}

// Update will check for new scores and update the ranking.
// Will not fetch from the API if it's been less than 30s since last fetch.
func (c *Client) Update() bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Rate limit update requests to 30 seconds
	if time.Now().Sub(c.LastUpdate) < 30*time.Second {
		return false
	}

	var scores = c.getRecentScores()
	c.log("Score count: ", len(scores))

	// Check if there are new scores
	if len(scores) == 0 || scores[0].CreatedAt.Compare(c.LastUpdate) <= 0 {
		c.log("No updates...")
		return false
	}

	// Add new scores to the ranks
	for _, score := range scores {
		if score.CreatedAt.Before(c.LastUpdate) {
			break
		}

		c.log("Potential new score: ", score.ID)
		c.ranking.AddScore(score)
	}

	// Update last signal
	c.LastUpdate = scores[0].CreatedAt
	return true
}

// Ranking safely returns client ranking.
// Modifications in the resulting ranking will not affect client.
func (c *Client) Ranking() player.Ranking {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.ranking
}

func (c *Client) log(v ...any) {
	var l = "PhantomClient." + c.Username + " : " + fmt.Sprint(v...)
	log.Println(l)
}

func (c *Client) getRecentScores() player.Scores {
	return osu.GetRecentScores(c.Provider.GetToken(), c.UserID)
}

func (c *Client) getBeatmapScores(beatmapID int) player.Scores {
	return osu.GetBeatmapScores(c.Provider.GetToken(), c.UserID, beatmapID)
}
