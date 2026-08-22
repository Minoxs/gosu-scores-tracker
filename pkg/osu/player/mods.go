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
