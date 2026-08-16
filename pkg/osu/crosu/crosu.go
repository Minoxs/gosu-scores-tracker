package crosu

/*
#cgo CFLAGS: -I ../../../deps/crosu-pp/bindings
#cgo LDFLAGS: -L ../../../deps/crosu-pp/target/release -lcrosu_pp
#include "crosu.h"
*/
import "C"
import (
	"strings"
	"unsafe"
)

type (
	GameMode C.GameMode
	ModType  C.u32
)

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

func (m ModType) ToC() C.u32 {
	return (C.u32)(m)
}

func (m GameMode) ToC() C.GameMode {
	return (C.GameMode)(m)
}

const (
	NF ModType = 1 << 0
	EZ         = 1 << 1
	TD         = 1 << 2
	HD         = 1 << 3
	HR         = 1 << 4
	DT         = 1 << 5
	RX         = 1 << 6
	HT         = 1 << 7
	FL         = 1 << 8
	SO         = 1 << 9
)

func (m ModType) FromString(s string) ModType {
	switch s {
	case "NF":
		return NF
	case "EZ":
		return EZ
	case "TD":
		return TD
	case "HD":
		return HD
	case "HR":
		return HR
	case "DT":
		return DT
	case "RX":
		return RX
	case "HT":
		return HT
	case "FL":
		return FL
	case "SO":
		return SO
	default:
		return 0
	}
}

func ModTypeFromStringArray(arr []string) (res ModType) {
	res = 0

	for _, s := range arr {
		res |= res.FromString(s)
	}

	return
}

// OsuDifficultyAttributes mirrors the osu!standard difficulty attributes
// returned by rosu-pp. Storing these lets pp be recomputed without the map file.
type OsuDifficultyAttributes struct {
	Aim                          float64
	AimDifficultSliderCount      float64
	Speed                        float64
	Flashlight                   float64
	SliderFactor                 float64
	AimTopWeightedSliderFactor   float64
	SpeedTopWeightedSliderFactor float64
	SpeedNoteCount               float64
	AimDifficultStrainCount      float64
	SpeedDifficultStrainCount    float64
	NestedScorePerObject         float64
	LegacyScoreBaseMultiplier    float64
	MaximumLegacyComboScore      float64
	AR                           float64
	GreatHitWindow               float64
	OkHitWindow                  float64
	MehHitWindow                 float64
	HP                           float64
	NCircles                     int
	NSliders                     int
	NLargeTicks                  int
	NSpinners                    int
	Stars                        float64
	MaxCombo                     int
}

func osuAttrFromC(c C.OsuDifficultyAttributes) OsuDifficultyAttributes {
	return OsuDifficultyAttributes{
		Aim:                          float64(c.aim),
		AimDifficultSliderCount:      float64(c.aimDifficultSliderCount),
		Speed:                        float64(c.speed),
		Flashlight:                   float64(c.flashlight),
		SliderFactor:                 float64(c.sliderFactor),
		AimTopWeightedSliderFactor:   float64(c.aimTopWeightedSliderFactor),
		SpeedTopWeightedSliderFactor: float64(c.speedTopWeightedSliderFactor),
		SpeedNoteCount:               float64(c.speedNoteCount),
		AimDifficultStrainCount:      float64(c.aimDifficultStrainCount),
		SpeedDifficultStrainCount:    float64(c.speedDifficultStrainCount),
		NestedScorePerObject:         float64(c.nestedScorePerObject),
		LegacyScoreBaseMultiplier:    float64(c.legacyScoreBaseMultiplier),
		MaximumLegacyComboScore:      float64(c.maximumLegacyComboScore),
		AR:                           float64(c.ar),
		GreatHitWindow:               float64(c.greatHitWindow),
		OkHitWindow:                  float64(c.okHitWindow),
		MehHitWindow:                 float64(c.mehHitWindow),
		HP:                           float64(c.hp),
		NCircles:                     int(c.nCircles),
		NSliders:                     int(c.nSliders),
		NLargeTicks:                  int(c.nLargeTicks),
		NSpinners:                    int(c.nSpinners),
		Stars:                        float64(c.stars),
		MaxCombo:                     int(c.maxCombo),
	}
}

func (a OsuDifficultyAttributes) toC() C.OsuDifficultyAttributes {
	return C.OsuDifficultyAttributes{
		aim:                          C.double(a.Aim),
		aimDifficultSliderCount:      C.double(a.AimDifficultSliderCount),
		speed:                        C.double(a.Speed),
		flashlight:                   C.double(a.Flashlight),
		sliderFactor:                 C.double(a.SliderFactor),
		aimTopWeightedSliderFactor:   C.double(a.AimTopWeightedSliderFactor),
		speedTopWeightedSliderFactor: C.double(a.SpeedTopWeightedSliderFactor),
		speedNoteCount:               C.double(a.SpeedNoteCount),
		aimDifficultStrainCount:      C.double(a.AimDifficultStrainCount),
		speedDifficultStrainCount:    C.double(a.SpeedDifficultStrainCount),
		nestedScorePerObject:         C.double(a.NestedScorePerObject),
		legacyScoreBaseMultiplier:    C.double(a.LegacyScoreBaseMultiplier),
		maximumLegacyComboScore:      C.double(a.MaximumLegacyComboScore),
		ar:                           C.double(a.AR),
		greatHitWindow:               C.double(a.GreatHitWindow),
		okHitWindow:                  C.double(a.OkHitWindow),
		mehHitWindow:                 C.double(a.MehHitWindow),
		hp:                           C.double(a.HP),
		nCircles:                     C.size_t(a.NCircles),
		nSliders:                     C.size_t(a.NSliders),
		nLargeTicks:                  C.size_t(a.NLargeTicks),
		nSpinners:                    C.size_t(a.NSpinners),
		stars:                        C.double(a.Stars),
		maxCombo:                     C.size_t(a.MaxCombo),
	}
}

// GetOsuDifficultyAttributes calculates osu!standard difficulty attributes for the given mods.
// Returns ok=false if difficulty calculation failed.
func GetOsuDifficultyAttributes(beatmap []byte, mods ModType) (attr OsuDifficultyAttributes, ok bool) {
	if beatmap[len(beatmap)-1] != 0 {
		beatmap = append(beatmap, 0)
	}

	var ptr = (*C.char)(unsafe.Pointer(&beatmap[0]))
	var res = C.GetOsuDifficultyAttributes(ptr, false, mods.ToC())
	if !bool(res.success) {
		return OsuDifficultyAttributes{}, false
	}

	return osuAttrFromC(res.attr), true
}

// GetOsuPP computes pp from cached osu!standard difficulty attributes.
// mods must match the ones used to calculate the attributes.
func GetOsuPP(attr OsuDifficultyAttributes, mods ModType, combo, n300, n100, n050, misses int) float64 {
	var score = C.Score{
		combo:  (C.size_t)(combo),
		n300:   (C.size_t)(n300),
		n100:   (C.size_t)(n100),
		n050:   (C.size_t)(n050),
		misses: (C.size_t)(misses),
	}
	return (float64)(C.GetOsuPP(attr.toC(), mods.ToC(), score))
}

func GetPPFromMap(beatmap []byte, combo, n300, n100, n050, misses int, mods ModType, mode GameMode) float64 {
	if beatmap[len(beatmap)-1] != 0 {
		beatmap = append(beatmap, 0)
	}

	var ptr = (*C.char)(unsafe.Pointer(&beatmap[0]))
	var score = C.Score{
		combo:  (C.size_t)(combo),
		n300:   (C.size_t)(n300),
		n100:   (C.size_t)(n100),
		n050:   (C.size_t)(n050),
		misses: (C.size_t)(misses),
	}
	return (float64)(C.GetPPFromMap(ptr, false, mods.ToC(), mode.ToC(), score))
}

func GetPPFromFile(filename string, combo, n300, n100, n050, misses int, mods ModType, mode GameMode) float64 {
	var beatmap = append([]byte(filename), 0)
	var ptr = (*C.char)(unsafe.Pointer(&beatmap[0]))
	var score = C.Score{
		combo:  (C.size_t)(combo),
		n300:   (C.size_t)(n300),
		n100:   (C.size_t)(n100),
		n050:   (C.size_t)(n050),
		misses: (C.size_t)(misses),
	}
	return (float64)(C.GetPPFromMap(ptr, true, mods.ToC(), mode.ToC(), score))
}
