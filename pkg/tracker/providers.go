package tracker

import "github.com/minoxs/gosu-api/pkg/gosu"

// ScoreProvider is a source of lean scores from the global feed, carrying only a
// beatmap_id. Drain a returned channel or delivery blocks; close it to unsubscribe.
// RealtimePoller implements it
type ScoreProvider interface {
	Subscribe() chan gosu.Score
}

// FullScoreProvider is a source of full scores with the beatmap and beatmapset
// embedded. Drain a returned channel or delivery blocks; close it to unsubscribe.
// UserPoller implements it
type FullScoreProvider interface {
	Subscribe() chan gosu.FullScore
}

var (
	_ ScoreProvider     = (*RealtimePoller)(nil)
	_ FullScoreProvider = (*UserPoller)(nil)
)
