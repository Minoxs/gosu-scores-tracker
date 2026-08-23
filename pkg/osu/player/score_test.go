package player

import (
	"encoding/json"
	"testing"
	"time"
)

// sample is a trimmed osu! API v2 solo_score object, keeping the fields Score
// decodes and their real shapes (mods as objects, lazer statistics, ended_at,
// ruleset_id, nested beatmap/beatmapset).
const sample = `{
  "id": 7027024157,
  "user_id": 30692023,
  "ended_at": "2026-07-07T03:30:08Z",
  "accuracy": 0.91954,
  "mods": [{"acronym": "HD"}, {"acronym": "DT"}],
  "total_score": 700980,
  "max_combo": 95,
  "rank": "A",
  "passed": true,
  "pp": 24.5157,
  "ruleset_id": 0,
  "statistics": {"great": 305, "ok": 30, "meh": 2, "miss": 13, "ignore_hit": 1},
  "beatmap": {"id": 2335023, "status": "ranked", "difficulty_rating": 6.2, "version": "Extra", "mode": "osu"},
  "beatmapset": {"id": 55, "title": "Song", "artist": "Artist", "creator": "Mapper"}
}`

func TestScoreDecode(t *testing.T) {
	var s Score
	if err := json.Unmarshal([]byte(sample), &s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if s.ID != 7027024157 || s.UserID != 30692023 {
		t.Errorf("id/user = %d/%d", s.ID, s.UserID)
	}
	if want := time.Date(2026, 7, 7, 3, 30, 8, 0, time.UTC); !s.EndedAt.Equal(want) {
		t.Errorf("EndedAt = %v, want %v", s.EndedAt, want)
	}
	if got := s.Mods.Acronyms(); len(got) != 2 || got[0] != "HD" || got[1] != "DT" {
		t.Errorf("Mods = %v, want [HD DT]", got)
	}
	if s.Statistics.Great != 305 || s.Statistics.Ok != 30 || s.Statistics.Meh != 2 || s.Statistics.Miss != 13 {
		t.Errorf("statistics = %+v, want 305/30/2/13", s.Statistics)
	}
	if s.TotalScore != 700980 {
		t.Errorf("TotalScore = %d, want 700980", s.TotalScore)
	}
	if s.Mode() != "osu" {
		t.Errorf("Mode() = %q, want osu (from ruleset_id 0)", s.Mode())
	}
	if s.PP != 24.5157 || !s.Passed {
		t.Errorf("pp/passed = %v/%v", s.PP, s.Passed)
	}
	if s.Beatmap.Status != StatusRanked || !s.AwardsPP() {
		t.Errorf("Beatmap.Status = %q, AwardsPP = %v", s.Beatmap.Status, s.AwardsPP())
	}
	if s.Beatmap.ID != 2335023 || s.BeatmapSet.Title != "Song" || s.BeatmapSet.Creator != "Mapper" {
		t.Errorf("beatmap/set = %d/%q/%q", s.Beatmap.ID, s.BeatmapSet.Title, s.BeatmapSet.Creator)
	}
}

func TestScore_AwardsPP(t *testing.T) {
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
		if got := score.AwardsPP(); got != want {
			t.Errorf("status %q: AwardsPP() = %v, want %v", status, got, want)
		}
	}
}

// A lean feed score carries the ranked flag but no beatmap, so AwardsPP must honor
// the flag before enrichment fills the status.
func TestScore_AwardsPP_LeanRankedFlag(t *testing.T) {
	lean := Score{Ranked: true}
	if !lean.AwardsPP() {
		t.Error("AwardsPP() = false for a ranked-flagged score with no beatmap, want true")
	}
}

// An unranked mod disqualifies pp even on an otherwise pp-awarding ranked map.
func TestScore_AwardsPP_UnrankedMods(t *testing.T) {
	for _, acronym := range []string{"RX", "AP", "DA", "SO"} {
		score := Score{Ranked: true, Beatmap: Beatmap{Status: StatusRanked}, Mods: Mods{{Acronym: acronym}}}
		if score.AwardsPP() {
			t.Errorf("AwardsPP() = true with unranked mod %q, want false", acronym)
		}
	}

	ranked := Score{Ranked: true, Beatmap: Beatmap{Status: StatusRanked}, Mods: Mods{{Acronym: "HD"}, {Acronym: "DT"}}}
	if !ranked.AwardsPP() {
		t.Error("AwardsPP() = false with ranked mods HD,DT, want true")
	}
}

func TestScoreMode(t *testing.T) {
	cases := map[int]string{0: "osu", 1: "taiko", 2: "fruits", 3: "mania", 9: ""}
	for id, want := range cases {
		if got := (Score{RulesetID: id}).Mode(); got != want {
			t.Errorf("ruleset %d Mode() = %q, want %q", id, got, want)
		}
	}
}
