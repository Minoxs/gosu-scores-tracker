package osu

type BeatmapMod string

const (
	// Ranked mods
	EZ BeatmapMod = "EZ"
	NF            = "NF"
	HT            = "HT"
	HR            = "HR"
	SD            = "SD"
	PF            = "PF"
	DT            = "DT"
	NC            = "NC"
	HD            = "HD"
	FL            = "FL"
	// Unranked mods
	RL  = "RL"
	AP  = "AP"
	SO  = "SO"
	AT  = "AT"
	CM  = "CM"
	SV2 = "SV2"
	TP  = "TP"
)

func (mod BeatmapMod) Ranked() bool {
	rankedMods := []BeatmapMod{EZ, NF, HT, HR, SD, PF, DT, NC, HD, FL}
	for _, m := range rankedMods {
		if m == mod {
			return true
		}
	}
	return false
}

func (mod BeatmapMod) CSModifier(cs float64) float64 {
	switch mod {
	case EZ:
		return cs / 2
	case HR:
		return cs * 1.3
	default:
		return cs
	}
}

func containsMod(arr []string, e BeatmapMod) bool {
	for _, x := range arr {
		if BeatmapMod(x) == e {
			return true
		}
	}
	return false
}
