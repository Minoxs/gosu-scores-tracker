package crosu

import "strings"

const (
	OSU   GameMode = 0
	TAIKO          = 1
	CATCH          = 2
	MANIA          = 3
)

func GameModeFromString(s string) GameMode {
	switch strings.ToLower(s) {
	case "osu":
		return OSU
	case "taiko":
		return TAIKO
	case "catch":
		return CATCH
	case "mania":
		return MANIA
	default:
		return OSU
	}
}
