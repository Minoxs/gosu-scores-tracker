package phantom

import (
	"errors"
	"github.com/minoxs/osu-phantom/pkg/osu"
	"github.com/minoxs/osu-phantom/pkg/osu/player"
	"log/slog"
	"sync"
	"time"
)

// TODO CONFIGURE START DATE
// TODO CACHE DIFFICULTY SCORE INSTEAD OF BEATMAP

type (
	// AuthProvider is required for requests which require OAuth authorization
	AuthProvider interface {
		GetToken() *osu.GuestToken
	}

	// NewScore contains information of a new score
	NewScore struct {
		Position int
		player.Score
	}

	// Client handles the tracking of a user's scores
	Client struct {
		UserID   int
		Username string
		Provider AuthProvider
		Logger   *slog.Logger

		OnNewScores func([]NewScore)
		ranking     player.Ranking
		LastUpdate  time.Time

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
	client.Logger = slog.Default().With("Username", username)

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
func (c *Client) KeepUpdated(checkInterval time.Duration, maxIdle time.Duration) {
	c.Logger.Info("Running KeepUpdated")
	defer func() {
		if r := recover(); r != nil {
			c.Logger.Error("Recovered from panic", "Panic", r)
		}
		c.Logger.Info("Stopping KeepUpdated routine")
	}()

	var interval = time.NewTimer(0)
	defer interval.Stop()

	for {
		select {
		case <-interval.C:
			c.Logger.Debug("Getting new scores")

			// Stop if idle for too long
			if !c.Update() && time.Now().Sub(c.LastUpdate) > maxIdle {
				return
			}

			c.Logger.Debug("New scores", "Ranking", c.ranking)

			// Reset timer
			interval.Reset(checkInterval)
		}
	}
}

// Update will check for new scores and update the ranking.
// Will not fetch from the API if it's been less than 30s since last fetch.
// Returns true when new scores were received from the API, even if they don't go into the ranking.
func (c *Client) Update() bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Rate limit update requests to 30 seconds
	if time.Now().Sub(c.LastUpdate) < 30*time.Second {
		return false
	}

	var scores = c.getRecentScores()
	c.Logger.Debug("Recent scores", "Count", len(scores))

	// Check if there are new scores
	if len(scores) == 0 || scores[0].CreatedAt.Compare(c.LastUpdate) <= 0 {
		c.Logger.Debug("No new scores")
		return false
	}
	c.processNewScores(scores)

	return true
}

// Ranking safely returns client ranking.
// Modifications in the resulting ranking will not affect client.
func (c *Client) Ranking() player.Ranking {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.ranking
}

func (c *Client) processNewScores(scores player.Scores) {
	// Keep track of added scores
	var (
		count     int
		newScores [10]NewScore
	)

	// Add new scores to the ranks
	for _, score := range scores {
		if score.CreatedAt.Compare(c.LastUpdate) <= 0 {
			break
		}

		osu.GetPP(&score)
		c.Logger.Debug("Possible new score", "ID", score.ID, "BeatmapID", score.Beatmap.ID, "Title", score.BeatmapSet.Title, "PP", score.PP)
		if rank, added := c.ranking.AddScore(score); added {
			c.Logger.Info("New score", "ID", score.ID, "BeatmapID", score.Beatmap.ID, "Title", score.BeatmapSet.Title, "PP", score.PP)
			newScores[count] = NewScore{
				Position: rank,
				Score:    score,
			}
			count++
		}
	}

	// Fire new scores event
	if c.OnNewScores != nil {
		c.OnNewScores(newScores[:count])
	}

	// Update last signal
	c.LastUpdate = scores[0].CreatedAt
}

func (c *Client) getRecentScores() player.Scores {
	return osu.GetRecentScores(c.Provider.GetToken(), c.UserID)
}

func (c *Client) getBeatmapScores(beatmapID int) player.Scores {
	return osu.GetBeatmapScores(c.Provider.GetToken(), c.UserID, beatmapID)
}
