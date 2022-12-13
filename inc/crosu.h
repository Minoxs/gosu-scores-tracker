#ifndef CROSU_PP
#define CROSU_PP

#include "stdint.h"
#include "stdbool.h"

typedef float f32;
typedef double f64;
typedef uint32_t u32;
typedef size_t usize;

const u32
	NF = 0b1,
	EZ = 0b10,
	TD = 0b100,
	HD = 0b1000,
	HR = 0b10000,
	DT = 0b100000,
	RX = 0b1000000,
	HT = 0b10000000,
	FL = 0b100000000,
	SO = 0b1000000000;

typedef struct {
	f64 aim;
	f64 speed;
	f64 flashlight;
	f64 sliderFactor;
	f64 speedNodeCount;
	f64 ar;
	f64 od;
	f64 hp;
	usize nCircles;
	usize nSliders;
	usize nSpinners;
	f64 stars;
	usize maxCombo;
} OsuDifficultyAttributes;

typedef struct {
	bool success;
	f64 pp;
	OsuDifficultyAttributes attr;
} OsuDiffResult;

typedef struct {
	usize combo;
	usize n300;
	usize n100;
	usize n050;
} Score;

typedef enum {
	Osu = 0,
	Taiko = 1,
	Catch = 2,
	Mania = 3
} GameMode;

// Functions of specific modes will return the required attributes
// to calculate pp without requiring beatmap file

// Osu
OsuDiffResult GetOsuDifficultyAttributes(char* ptr, bool isFile, u32 mods);
f64 GetOsuPP(OsuDifficultyAttributes diff, Score score);
OsuDiffResult GetOsuPPFromMap(char* ptr, bool isFile, u32 mods, Score score);

// Generic
// Returns PP, but requires recalculating map difficulty every time
f64 GetPPFromMap(char* ptr, bool isFile, u32 mods, GameMode mode, Score score);

#endif // CROSU_PP
