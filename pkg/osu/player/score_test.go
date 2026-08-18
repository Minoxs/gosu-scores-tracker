package player

import "testing"

func TestScore_IsRanked(t *testing.T) {
	cases := map[BeatmapStatus]bool{
		"ranked":    true,
		"approved":  true,
		"loved":     false,
		"qualified": false,
		"pending":   false,
		"wip":       false,
		"graveyard": false,
		"":          false,
	}

	for status, want := range cases {
		score := Score{Beatmap: Beatmap{Status: status}}
		if got := score.IsRanked(); got != want {
			t.Errorf("status %q: IsRanked() = %v, want %v", status, got, want)
		}
	}
}
