package player

// Mod is one mod on a score, an object in the osu! API v2 shape.
type Mod struct {
	Acronym string `json:"acronym"`
}

// Mods is a score's mod list in the osu! API v2 shape, an array of objects.
type Mods []Mod

// Acronyms returns the mod acronyms in order, e.g. ["HD","DT"].
func (m Mods) Acronyms() []string {
	out := make([]string, len(m))
	for i := range m {
		out[i] = m[i].Acronym
	}
	return out
}

// unrankedMods are the osu! mods that make a score earn no pp, so a play using any
// of them is disqualified even on a ranked map. osu! reports pp of 0 for these, so
// without this gate a crosu fallback would invent a nonzero pp for them. This is a
// denylist: extend it when osu! ships another unranked mod.
var unrankedMods = map[string]bool{
	"RX": true, // Relax
	"AP": true, // Autopilot
	"SO": true, // Spun Out
	"TP": true, // Target Practice
	"DA": true, // Difficulty Adjust
	"RD": true, // Random
	"MG": true, // Magnetised
	"RP": true, // Repel
	"AS": true, // Adaptive Speed
	"FR": true, // Freeze Frame
	"BU": true, // Bubbles
	"SY": true, // Synesthesia
	"DP": true, // Depth
	"TR": true, // Transform
	"WG": true, // Wiggle
	"BR": true, // Barrel Roll
	"WU": true, // Wind Up
	"WD": true, // Wind Down
}

// HasUnranked reports whether any mod in the list disqualifies the score from pp.
func (m Mods) HasUnranked() bool {
	for i := range m {
		if unrankedMods[m[i].Acronym] {
			return true
		}
	}
	return false
}
