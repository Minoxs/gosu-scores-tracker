package phantom

import (
	"sync"

	"github.com/minoxs/osu-phantom/pkg/osu/player"
)

// ScoreProvider is a source of score events. It surfaces scores without regard
// to who is tracked: kay's oSC feed surfaces every global score, an osu-API
// provider surfaces a user's scores when polled. Each Subscribe call returns an
// independent stream, so several consumers can read the same feed at once, for
// example one folding scores into rankings while another auto-tracks new users.
type ScoreProvider interface {
	Subscribe() <-chan player.Score
}

// DefaultSubBuffer is the per-subscriber channel capacity a Broadcaster uses.
const DefaultSubBuffer = 256

// Broadcaster is a reusable ScoreProvider. A concrete provider embeds one and
// calls Emit for each score it obtains; the Broadcaster fans that score out to
// every current subscriber. Delivery is best-effort: a subscriber that does not
// keep draining its channel misses events once its buffer fills, rather than
// stalling the feed for the other subscribers.
type Broadcaster struct {
	buffer int
	mu     sync.Mutex
	subs   map[chan player.Score]struct{}
}

// NewBroadcaster builds an empty Broadcaster with the default subscriber buffer.
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		buffer: DefaultSubBuffer,
		subs:   make(map[chan player.Score]struct{}),
	}
}

// Subscribe registers a new consumer and returns its stream. Close ends every
// stream a Broadcaster has handed out.
func (b *Broadcaster) Subscribe() <-chan player.Score {
	ch := make(chan player.Score, b.buffer)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

// Emit fans a score out to every subscriber, skipping any whose buffer is full.
func (b *Broadcaster) Emit(s player.Score) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.subs {
		select {
		case ch <- s:
		default:
		}
	}
}

// Close drops and closes every subscription so their readers finish. The
// Broadcaster may still be subscribed to again afterwards.
func (b *Broadcaster) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.subs {
		close(ch)
		delete(b.subs, ch)
	}
}
