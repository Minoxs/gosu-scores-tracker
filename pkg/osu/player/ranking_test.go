package player

import (
	"log"
	"testing"
)

type TestCaseInfo struct {
	Name          string
	Scores        Scores
	ExpectedCount int
}

func assertEqual[T comparable](t *testing.T, val1 T, val2 T) {
	if val1 != val2 {
		t.Fatalf("Values are not equal\nExpected=%v\nActual=%v\n", val1, val2)
	}
}

func assert(t *testing.T, test bool, msg string) {
	if !test {
		t.Fatal(msg)
	}
}

func getTestCases() []TestCaseInfo {
	return []TestCaseInfo{
		{
			Name: "Single",
			Scores: Scores{
				{ID: 1, Beatmap: Beatmap{ID: 1}, PP: 250.0},
			},
			ExpectedCount: 1,
		},
		{
			Name: "Triple",
			Scores: Scores{
				{ID: 1, Beatmap: Beatmap{ID: 1}, PP: 369.15}, {ID: 2, Beatmap: Beatmap{ID: 2}, PP: 405.17}, {ID: 3, Beatmap: Beatmap{ID: 3}, PP: 5.89},
			},
			ExpectedCount: 3,
		},
		{
			Name: "Six scores",
			Scores: Scores{
				{ID: 1, Beatmap: Beatmap{ID: 1}, PP: 690.64}, {ID: 2, Beatmap: Beatmap{ID: 2}, PP: 216.06}, {ID: 3, Beatmap: Beatmap{ID: 3}, PP: 141.08},
				{ID: 4, Beatmap: Beatmap{ID: 4}, PP: 6.03}, {ID: 5, Beatmap: Beatmap{ID: 5}, PP: 14.58}, {ID: 6, Beatmap: Beatmap{ID: 6}, PP: 390.45},
			},
			ExpectedCount: 6,
		},
		{
			Name: "Big ranking",
			Scores: Scores{
				{ID: 1, Beatmap: Beatmap{ID: 1}, PP: 81.98}, {ID: 2, Beatmap: Beatmap{ID: 2}, PP: 106.62}, {ID: 3, Beatmap: Beatmap{ID: 3}, PP: 105.02}, {ID: 4, Beatmap: Beatmap{ID: 4}, PP: 70.37},
				{ID: 5, Beatmap: Beatmap{ID: 5}, PP: 176.07}, {ID: 6, Beatmap: Beatmap{ID: 6}, PP: 278.56}, {ID: 7, Beatmap: Beatmap{ID: 7}, PP: 109.50}, {ID: 8, Beatmap: Beatmap{ID: 8}, PP: 211.35},
				{ID: 9, Beatmap: Beatmap{ID: 9}, PP: 55.08}, {ID: 10, Beatmap: Beatmap{ID: 10}, PP: 22.45}, {ID: 11, Beatmap: Beatmap{ID: 11}, PP: 132.05}, {ID: 12, Beatmap: Beatmap{ID: 12}, PP: 129.05},
				{ID: 13, Beatmap: Beatmap{ID: 13}, PP: 826.40}, {ID: 14, Beatmap: Beatmap{ID: 14}, PP: 163.30}, {ID: 15, Beatmap: Beatmap{ID: 15}, PP: 57.45}, {ID: 16, Beatmap: Beatmap{ID: 16}, PP: 13.82},
				{ID: 17, Beatmap: Beatmap{ID: 17}, PP: 339.63}, {ID: 18, Beatmap: Beatmap{ID: 18}, PP: 584.32}, {ID: 19, Beatmap: Beatmap{ID: 19}, PP: 477.58}, {ID: 20, Beatmap: Beatmap{ID: 20}, PP: 204.79},
				{ID: 21, Beatmap: Beatmap{ID: 21}, PP: 95.32}, {ID: 22, Beatmap: Beatmap{ID: 22}, PP: 168.66}, {ID: 23, Beatmap: Beatmap{ID: 23}, PP: 376.65}, {ID: 24, Beatmap: Beatmap{ID: 24}, PP: 36.46},
				{ID: 25, Beatmap: Beatmap{ID: 25}, PP: 162.12}, {ID: 26, Beatmap: Beatmap{ID: 26}, PP: 504.60}, {ID: 27, Beatmap: Beatmap{ID: 27}, PP: 53.14}, {ID: 28, Beatmap: Beatmap{ID: 28}, PP: 684.83},
			},
			ExpectedCount: 28,
		},
		{
			Name: "Exactly 100",
			Scores: Scores{
				{ID: 1, Beatmap: Beatmap{ID: 1}, PP: 1.91}, {ID: 2, Beatmap: Beatmap{ID: 2}, PP: 432.57}, {ID: 3, Beatmap: Beatmap{ID: 3}, PP: 777.75}, {ID: 4, Beatmap: Beatmap{ID: 4}, PP: 227.26},
				{ID: 5, Beatmap: Beatmap{ID: 5}, PP: 377.28}, {ID: 6, Beatmap: Beatmap{ID: 6}, PP: 127.23}, {ID: 7, Beatmap: Beatmap{ID: 7}, PP: 223.17}, {ID: 8, Beatmap: Beatmap{ID: 8}, PP: 543.68},
				{ID: 9, Beatmap: Beatmap{ID: 9}, PP: 22.46}, {ID: 10, Beatmap: Beatmap{ID: 10}, PP: 809.20}, {ID: 11, Beatmap: Beatmap{ID: 11}, PP: 1.50}, {ID: 12, Beatmap: Beatmap{ID: 12}, PP: 86.67},
				{ID: 13, Beatmap: Beatmap{ID: 13}, PP: 190.85}, {ID: 14, Beatmap: Beatmap{ID: 14}, PP: 151.16}, {ID: 15, Beatmap: Beatmap{ID: 15}, PP: 669.64}, {ID: 16, Beatmap: Beatmap{ID: 16}, PP: 176.70},
				{ID: 17, Beatmap: Beatmap{ID: 17}, PP: 11.91}, {ID: 18, Beatmap: Beatmap{ID: 18}, PP: 56.34}, {ID: 19, Beatmap: Beatmap{ID: 19}, PP: 299.09}, {ID: 20, Beatmap: Beatmap{ID: 20}, PP: 92.60},
				{ID: 21, Beatmap: Beatmap{ID: 21}, PP: 174.75}, {ID: 22, Beatmap: Beatmap{ID: 22}, PP: 219.66}, {ID: 23, Beatmap: Beatmap{ID: 23}, PP: 6.62}, {ID: 24, Beatmap: Beatmap{ID: 24}, PP: 330.46},
				{ID: 25, Beatmap: Beatmap{ID: 25}, PP: 50.34}, {ID: 26, Beatmap: Beatmap{ID: 26}, PP: 132.53}, {ID: 27, Beatmap: Beatmap{ID: 27}, PP: 383.88}, {ID: 28, Beatmap: Beatmap{ID: 28}, PP: 2.65},
				{ID: 29, Beatmap: Beatmap{ID: 29}, PP: 2.68}, {ID: 30, Beatmap: Beatmap{ID: 30}, PP: 12.89}, {ID: 31, Beatmap: Beatmap{ID: 31}, PP: 31.87}, {ID: 32, Beatmap: Beatmap{ID: 32}, PP: 77.76},
				{ID: 33, Beatmap: Beatmap{ID: 33}, PP: 206.75}, {ID: 34, Beatmap: Beatmap{ID: 34}, PP: 662.69}, {ID: 35, Beatmap: Beatmap{ID: 35}, PP: 475.25}, {ID: 36, Beatmap: Beatmap{ID: 36}, PP: 207.83},
				{ID: 37, Beatmap: Beatmap{ID: 37}, PP: 250.62}, {ID: 38, Beatmap: Beatmap{ID: 38}, PP: 171.97}, {ID: 39, Beatmap: Beatmap{ID: 39}, PP: 108.74}, {ID: 40, Beatmap: Beatmap{ID: 40}, PP: 378.44},
				{ID: 41, Beatmap: Beatmap{ID: 41}, PP: 478.57}, {ID: 42, Beatmap: Beatmap{ID: 42}, PP: 260.13}, {ID: 43, Beatmap: Beatmap{ID: 43}, PP: 44.08}, {ID: 44, Beatmap: Beatmap{ID: 44}, PP: 30.77},
				{ID: 45, Beatmap: Beatmap{ID: 45}, PP: 415.29}, {ID: 46, Beatmap: Beatmap{ID: 46}, PP: 398.12}, {ID: 47, Beatmap: Beatmap{ID: 47}, PP: 241.62}, {ID: 48, Beatmap: Beatmap{ID: 48}, PP: 698.11},
				{ID: 49, Beatmap: Beatmap{ID: 49}, PP: 273.17}, {ID: 50, Beatmap: Beatmap{ID: 50}, PP: 47.67}, {ID: 51, Beatmap: Beatmap{ID: 51}, PP: 345.72}, {ID: 52, Beatmap: Beatmap{ID: 52}, PP: 81.74},
				{ID: 53, Beatmap: Beatmap{ID: 53}, PP: 739.04}, {ID: 54, Beatmap: Beatmap{ID: 54}, PP: 624.40}, {ID: 55, Beatmap: Beatmap{ID: 55}, PP: 652.44}, {ID: 56, Beatmap: Beatmap{ID: 56}, PP: 533.48},
				{ID: 57, Beatmap: Beatmap{ID: 57}, PP: 818.32}, {ID: 58, Beatmap: Beatmap{ID: 58}, PP: 13.05}, {ID: 59, Beatmap: Beatmap{ID: 59}, PP: 125.27}, {ID: 60, Beatmap: Beatmap{ID: 60}, PP: 456.04},
				{ID: 61, Beatmap: Beatmap{ID: 61}, PP: 386.51}, {ID: 62, Beatmap: Beatmap{ID: 62}, PP: 56.35}, {ID: 63, Beatmap: Beatmap{ID: 63}, PP: 272.82}, {ID: 64, Beatmap: Beatmap{ID: 64}, PP: 874.52},
				{ID: 65, Beatmap: Beatmap{ID: 65}, PP: 263.19}, {ID: 66, Beatmap: Beatmap{ID: 66}, PP: 84.93}, {ID: 67, Beatmap: Beatmap{ID: 67}, PP: 17.17}, {ID: 68, Beatmap: Beatmap{ID: 68}, PP: 630.61},
				{ID: 69, Beatmap: Beatmap{ID: 69}, PP: 222.87}, {ID: 70, Beatmap: Beatmap{ID: 70}, PP: 44.33}, {ID: 71, Beatmap: Beatmap{ID: 71}, PP: 555.45}, {ID: 72, Beatmap: Beatmap{ID: 72}, PP: 9.27},
				{ID: 73, Beatmap: Beatmap{ID: 73}, PP: 845.36}, {ID: 74, Beatmap: Beatmap{ID: 74}, PP: 239.58}, {ID: 75, Beatmap: Beatmap{ID: 75}, PP: 396.26}, {ID: 76, Beatmap: Beatmap{ID: 76}, PP: 48.83},
				{ID: 77, Beatmap: Beatmap{ID: 77}, PP: 348.94}, {ID: 78, Beatmap: Beatmap{ID: 78}, PP: 4.46}, {ID: 79, Beatmap: Beatmap{ID: 79}, PP: 889.00}, {ID: 80, Beatmap: Beatmap{ID: 80}, PP: 397.88},
				{ID: 81, Beatmap: Beatmap{ID: 81}, PP: 244.27}, {ID: 82, Beatmap: Beatmap{ID: 82}, PP: 202.41}, {ID: 83, Beatmap: Beatmap{ID: 83}, PP: 382.12}, {ID: 84, Beatmap: Beatmap{ID: 84}, PP: 639.83},
				{ID: 85, Beatmap: Beatmap{ID: 85}, PP: 58.52}, {ID: 86, Beatmap: Beatmap{ID: 86}, PP: 656.98}, {ID: 87, Beatmap: Beatmap{ID: 87}, PP: 769.99}, {ID: 88, Beatmap: Beatmap{ID: 88}, PP: 124.50},
				{ID: 89, Beatmap: Beatmap{ID: 89}, PP: 168.77}, {ID: 90, Beatmap: Beatmap{ID: 90}, PP: 719.15}, {ID: 91, Beatmap: Beatmap{ID: 91}, PP: 159.54}, {ID: 92, Beatmap: Beatmap{ID: 92}, PP: 405.20},
				{ID: 93, Beatmap: Beatmap{ID: 93}, PP: 659.02}, {ID: 94, Beatmap: Beatmap{ID: 94}, PP: 216.44}, {ID: 95, Beatmap: Beatmap{ID: 95}, PP: 107.84}, {ID: 96, Beatmap: Beatmap{ID: 96}, PP: 49.88},
				{ID: 97, Beatmap: Beatmap{ID: 97}, PP: 134.15}, {ID: 98, Beatmap: Beatmap{ID: 98}, PP: 498.19}, {ID: 99, Beatmap: Beatmap{ID: 99}, PP: 6.50}, {ID: 100, Beatmap: Beatmap{ID: 100}, PP: 281.48},
			},
			ExpectedCount: 100,
		},
		{
			Name: "More than a full ranking",
			Scores: Scores{
				{ID: 1, Beatmap: Beatmap{ID: 1}, PP: 295.32}, {ID: 2, Beatmap: Beatmap{ID: 2}, PP: 3.43}, {ID: 3, Beatmap: Beatmap{ID: 3}, PP: 67.31}, {ID: 4, Beatmap: Beatmap{ID: 4}, PP: 727.35},
				{ID: 5, Beatmap: Beatmap{ID: 5}, PP: 239.13}, {ID: 6, Beatmap: Beatmap{ID: 6}, PP: 73.64}, {ID: 7, Beatmap: Beatmap{ID: 7}, PP: 339.36}, {ID: 8, Beatmap: Beatmap{ID: 8}, PP: 506.01},
				{ID: 9, Beatmap: Beatmap{ID: 9}, PP: 275.33}, {ID: 10, Beatmap: Beatmap{ID: 10}, PP: 285.41}, {ID: 11, Beatmap: Beatmap{ID: 11}, PP: 53.52}, {ID: 12, Beatmap: Beatmap{ID: 12}, PP: 143.05},
				{ID: 13, Beatmap: Beatmap{ID: 13}, PP: 20.08}, {ID: 14, Beatmap: Beatmap{ID: 14}, PP: 80.64}, {ID: 15, Beatmap: Beatmap{ID: 15}, PP: 109.70}, {ID: 16, Beatmap: Beatmap{ID: 16}, PP: 10.06},
				{ID: 17, Beatmap: Beatmap{ID: 17}, PP: 89.00}, {ID: 18, Beatmap: Beatmap{ID: 18}, PP: 6.36}, {ID: 19, Beatmap: Beatmap{ID: 19}, PP: 33.29}, {ID: 20, Beatmap: Beatmap{ID: 20}, PP: 11.78},
				{ID: 21, Beatmap: Beatmap{ID: 21}, PP: 74.18}, {ID: 22, Beatmap: Beatmap{ID: 22}, PP: 19.23}, {ID: 23, Beatmap: Beatmap{ID: 23}, PP: 43.69}, {ID: 24, Beatmap: Beatmap{ID: 24}, PP: 312.58},
				{ID: 25, Beatmap: Beatmap{ID: 25}, PP: 107.72}, {ID: 26, Beatmap: Beatmap{ID: 26}, PP: 162.16}, {ID: 27, Beatmap: Beatmap{ID: 27}, PP: 61.40}, {ID: 28, Beatmap: Beatmap{ID: 28}, PP: 338.30},
				{ID: 29, Beatmap: Beatmap{ID: 29}, PP: 441.35}, {ID: 30, Beatmap: Beatmap{ID: 30}, PP: 285.01}, {ID: 31, Beatmap: Beatmap{ID: 31}, PP: 409.12}, {ID: 32, Beatmap: Beatmap{ID: 32}, PP: 160.62},
				{ID: 33, Beatmap: Beatmap{ID: 33}, PP: 863.42}, {ID: 34, Beatmap: Beatmap{ID: 34}, PP: 250.04}, {ID: 35, Beatmap: Beatmap{ID: 35}, PP: 157.38}, {ID: 36, Beatmap: Beatmap{ID: 36}, PP: 269.93},
				{ID: 37, Beatmap: Beatmap{ID: 37}, PP: 473.43}, {ID: 38, Beatmap: Beatmap{ID: 38}, PP: 511.64}, {ID: 39, Beatmap: Beatmap{ID: 39}, PP: 173.70}, {ID: 40, Beatmap: Beatmap{ID: 40}, PP: 762.89},
				{ID: 41, Beatmap: Beatmap{ID: 41}, PP: 93.87}, {ID: 42, Beatmap: Beatmap{ID: 42}, PP: 321.12}, {ID: 43, Beatmap: Beatmap{ID: 43}, PP: 78.76}, {ID: 44, Beatmap: Beatmap{ID: 44}, PP: 114.12},
				{ID: 45, Beatmap: Beatmap{ID: 45}, PP: 93.97}, {ID: 46, Beatmap: Beatmap{ID: 46}, PP: 60.55}, {ID: 47, Beatmap: Beatmap{ID: 47}, PP: 176.70}, {ID: 48, Beatmap: Beatmap{ID: 48}, PP: 98.16},
				{ID: 49, Beatmap: Beatmap{ID: 49}, PP: 128.36}, {ID: 50, Beatmap: Beatmap{ID: 50}, PP: 277.15}, {ID: 51, Beatmap: Beatmap{ID: 51}, PP: 380.59}, {ID: 52, Beatmap: Beatmap{ID: 52}, PP: 96.26},
				{ID: 53, Beatmap: Beatmap{ID: 53}, PP: 252.67}, {ID: 54, Beatmap: Beatmap{ID: 54}, PP: 106.80}, {ID: 55, Beatmap: Beatmap{ID: 55}, PP: 327.50}, {ID: 56, Beatmap: Beatmap{ID: 56}, PP: 238.93},
				{ID: 57, Beatmap: Beatmap{ID: 57}, PP: 503.40}, {ID: 58, Beatmap: Beatmap{ID: 58}, PP: 130.83}, {ID: 59, Beatmap: Beatmap{ID: 59}, PP: 72.63}, {ID: 60, Beatmap: Beatmap{ID: 60}, PP: 296.09},
				{ID: 61, Beatmap: Beatmap{ID: 61}, PP: 455.34}, {ID: 62, Beatmap: Beatmap{ID: 62}, PP: 347.66}, {ID: 63, Beatmap: Beatmap{ID: 63}, PP: 804.75}, {ID: 64, Beatmap: Beatmap{ID: 64}, PP: 14.87},
				{ID: 65, Beatmap: Beatmap{ID: 65}, PP: 680.41}, {ID: 66, Beatmap: Beatmap{ID: 66}, PP: 23.43}, {ID: 67, Beatmap: Beatmap{ID: 67}, PP: 567.78}, {ID: 68, Beatmap: Beatmap{ID: 68}, PP: 69.89},
				{ID: 69, Beatmap: Beatmap{ID: 69}, PP: 10.01}, {ID: 70, Beatmap: Beatmap{ID: 70}, PP: 531.28}, {ID: 71, Beatmap: Beatmap{ID: 71}, PP: 190.69}, {ID: 72, Beatmap: Beatmap{ID: 72}, PP: 84.50},
				{ID: 73, Beatmap: Beatmap{ID: 73}, PP: 64.58}, {ID: 74, Beatmap: Beatmap{ID: 74}, PP: 231.80}, {ID: 75, Beatmap: Beatmap{ID: 75}, PP: 35.36}, {ID: 76, Beatmap: Beatmap{ID: 76}, PP: 7.46},
				{ID: 77, Beatmap: Beatmap{ID: 77}, PP: 148.84}, {ID: 78, Beatmap: Beatmap{ID: 78}, PP: 159.55}, {ID: 79, Beatmap: Beatmap{ID: 79}, PP: 535.69}, {ID: 80, Beatmap: Beatmap{ID: 80}, PP: 218.39},
				{ID: 81, Beatmap: Beatmap{ID: 81}, PP: 648.71}, {ID: 82, Beatmap: Beatmap{ID: 82}, PP: 75.74}, {ID: 83, Beatmap: Beatmap{ID: 83}, PP: 224.10}, {ID: 84, Beatmap: Beatmap{ID: 84}, PP: 193.24},
				{ID: 85, Beatmap: Beatmap{ID: 85}, PP: 661.04}, {ID: 86, Beatmap: Beatmap{ID: 86}, PP: 8.46}, {ID: 87, Beatmap: Beatmap{ID: 87}, PP: 449.28}, {ID: 88, Beatmap: Beatmap{ID: 88}, PP: 883.39},
				{ID: 89, Beatmap: Beatmap{ID: 89}, PP: 345.51}, {ID: 90, Beatmap: Beatmap{ID: 90}, PP: 79.00}, {ID: 91, Beatmap: Beatmap{ID: 91}, PP: 128.69}, {ID: 92, Beatmap: Beatmap{ID: 92}, PP: 484.31},
				{ID: 93, Beatmap: Beatmap{ID: 93}, PP: 414.34}, {ID: 94, Beatmap: Beatmap{ID: 94}, PP: 66.90}, {ID: 95, Beatmap: Beatmap{ID: 95}, PP: 138.82}, {ID: 96, Beatmap: Beatmap{ID: 96}, PP: 266.95},
				{ID: 97, Beatmap: Beatmap{ID: 97}, PP: 773.66}, {ID: 98, Beatmap: Beatmap{ID: 98}, PP: 93.27}, {ID: 99, Beatmap: Beatmap{ID: 99}, PP: 128.80}, {ID: 100, Beatmap: Beatmap{ID: 100}, PP: 635.03},
				{ID: 101, Beatmap: Beatmap{ID: 101}, PP: 60.49}, {ID: 102, Beatmap: Beatmap{ID: 102}, PP: 28.98}, {ID: 103, Beatmap: Beatmap{ID: 103}, PP: 71.95}, {ID: 104, Beatmap: Beatmap{ID: 104}, PP: 343.90},
				{ID: 105, Beatmap: Beatmap{ID: 105}, PP: 70.33}, {ID: 106, Beatmap: Beatmap{ID: 106}, PP: 175.31}, {ID: 107, Beatmap: Beatmap{ID: 107}, PP: 220.89}, {ID: 108, Beatmap: Beatmap{ID: 108}, PP: 221.71},
				{ID: 109, Beatmap: Beatmap{ID: 109}, PP: 577.18}, {ID: 110, Beatmap: Beatmap{ID: 110}, PP: 36.72}, {ID: 111, Beatmap: Beatmap{ID: 111}, PP: 334.91}, {ID: 112, Beatmap: Beatmap{ID: 112}, PP: 49.95},
				{ID: 113, Beatmap: Beatmap{ID: 113}, PP: 7.87}, {ID: 114, Beatmap: Beatmap{ID: 114}, PP: 518.80}, {ID: 115, Beatmap: Beatmap{ID: 115}, PP: 424.40}, {ID: 116, Beatmap: Beatmap{ID: 116}, PP: 5.53},
				{ID: 117, Beatmap: Beatmap{ID: 117}, PP: 165.87}, {ID: 118, Beatmap: Beatmap{ID: 118}, PP: 66.41}, {ID: 119, Beatmap: Beatmap{ID: 119}, PP: 19.74}, {ID: 120, Beatmap: Beatmap{ID: 120}, PP: 457.63},
				{ID: 121, Beatmap: Beatmap{ID: 121}, PP: 71.28}, {ID: 122, Beatmap: Beatmap{ID: 122}, PP: 91.51}, {ID: 123, Beatmap: Beatmap{ID: 123}, PP: 257.60}, {ID: 124, Beatmap: Beatmap{ID: 124}, PP: 629.34},
				{ID: 125, Beatmap: Beatmap{ID: 125}, PP: 21.94}, {ID: 126, Beatmap: Beatmap{ID: 126}, PP: 190.83}, {ID: 127, Beatmap: Beatmap{ID: 127}, PP: 483.70}, {ID: 128, Beatmap: Beatmap{ID: 128}, PP: 458.06},
				{ID: 129, Beatmap: Beatmap{ID: 129}, PP: 215.51}, {ID: 130, Beatmap: Beatmap{ID: 130}, PP: 104.79}, {ID: 131, Beatmap: Beatmap{ID: 131}, PP: 118.05}, {ID: 132, Beatmap: Beatmap{ID: 132}, PP: 674.71},
				{ID: 133, Beatmap: Beatmap{ID: 133}, PP: 111.09}, {ID: 134, Beatmap: Beatmap{ID: 134}, PP: 35.06}, {ID: 135, Beatmap: Beatmap{ID: 135}, PP: 258.01}, {ID: 136, Beatmap: Beatmap{ID: 136}, PP: 154.23},
				{ID: 137, Beatmap: Beatmap{ID: 137}, PP: 693.71}, {ID: 138, Beatmap: Beatmap{ID: 138}, PP: 418.39}, {ID: 139, Beatmap: Beatmap{ID: 139}, PP: 304.57}, {ID: 140, Beatmap: Beatmap{ID: 140}, PP: 205.26},
				{ID: 141, Beatmap: Beatmap{ID: 141}, PP: 139.14}, {ID: 142, Beatmap: Beatmap{ID: 142}, PP: 37.05}, {ID: 143, Beatmap: Beatmap{ID: 143}, PP: 12.25}, {ID: 144, Beatmap: Beatmap{ID: 144}, PP: 457.93},
				{ID: 145, Beatmap: Beatmap{ID: 145}, PP: 696.30}, {ID: 146, Beatmap: Beatmap{ID: 146}, PP: 54.09}, {ID: 147, Beatmap: Beatmap{ID: 147}, PP: 411.07}, {ID: 148, Beatmap: Beatmap{ID: 148}, PP: 245.89},
				{ID: 149, Beatmap: Beatmap{ID: 149}, PP: 64.37}, {ID: 150, Beatmap: Beatmap{ID: 150}, PP: 102.83}, {ID: 151, Beatmap: Beatmap{ID: 151}, PP: 404.29}, {ID: 152, Beatmap: Beatmap{ID: 152}, PP: 11.07},
				{ID: 153, Beatmap: Beatmap{ID: 153}, PP: 340.52}, {ID: 154, Beatmap: Beatmap{ID: 154}, PP: 54.60}, {ID: 155, Beatmap: Beatmap{ID: 155}, PP: 137.16}, {ID: 156, Beatmap: Beatmap{ID: 156}, PP: 60.82},
				{ID: 157, Beatmap: Beatmap{ID: 157}, PP: 519.73}, {ID: 158, Beatmap: Beatmap{ID: 158}, PP: 6.32}, {ID: 159, Beatmap: Beatmap{ID: 159}, PP: 199.61}, {ID: 160, Beatmap: Beatmap{ID: 160}, PP: 180.80},
				{ID: 161, Beatmap: Beatmap{ID: 161}, PP: 708.30}, {ID: 162, Beatmap: Beatmap{ID: 162}, PP: 132.08}, {ID: 163, Beatmap: Beatmap{ID: 163}, PP: 207.74}, {ID: 164, Beatmap: Beatmap{ID: 164}, PP: 2.38},
				{ID: 165, Beatmap: Beatmap{ID: 165}, PP: 113.86}, {ID: 166, Beatmap: Beatmap{ID: 166}, PP: 99.39}, {ID: 167, Beatmap: Beatmap{ID: 167}, PP: 19.53}, {ID: 168, Beatmap: Beatmap{ID: 168}, PP: 120.19},
				{ID: 169, Beatmap: Beatmap{ID: 169}, PP: 446.42}, {ID: 170, Beatmap: Beatmap{ID: 170}, PP: 329.57}, {ID: 171, Beatmap: Beatmap{ID: 171}, PP: 163.00}, {ID: 172, Beatmap: Beatmap{ID: 172}, PP: 152.84},
				{ID: 173, Beatmap: Beatmap{ID: 173}, PP: 24.66}, {ID: 174, Beatmap: Beatmap{ID: 174}, PP: 88.35}, {ID: 175, Beatmap: Beatmap{ID: 175}, PP: 240.61}, {ID: 176, Beatmap: Beatmap{ID: 176}, PP: 206.67},
				{ID: 177, Beatmap: Beatmap{ID: 177}, PP: 553.56}, {ID: 178, Beatmap: Beatmap{ID: 178}, PP: 135.32}, {ID: 179, Beatmap: Beatmap{ID: 179}, PP: 326.01}, {ID: 180, Beatmap: Beatmap{ID: 180}, PP: 189.57},
				{ID: 181, Beatmap: Beatmap{ID: 181}, PP: 161.38}, {ID: 182, Beatmap: Beatmap{ID: 182}, PP: 56.32}, {ID: 183, Beatmap: Beatmap{ID: 183}, PP: 145.64}, {ID: 184, Beatmap: Beatmap{ID: 184}, PP: 81.48},
				{ID: 185, Beatmap: Beatmap{ID: 185}, PP: 49.07}, {ID: 186, Beatmap: Beatmap{ID: 186}, PP: 382.80}, {ID: 187, Beatmap: Beatmap{ID: 187}, PP: 255.71}, {ID: 188, Beatmap: Beatmap{ID: 188}, PP: 436.79},
				{ID: 189, Beatmap: Beatmap{ID: 189}, PP: 423.80}, {ID: 190, Beatmap: Beatmap{ID: 190}, PP: 106.30}, {ID: 191, Beatmap: Beatmap{ID: 191}, PP: 128.88}, {ID: 192, Beatmap: Beatmap{ID: 192}, PP: 2.01},
				{ID: 193, Beatmap: Beatmap{ID: 193}, PP: 45.14}, {ID: 194, Beatmap: Beatmap{ID: 194}, PP: 105.41}, {ID: 195, Beatmap: Beatmap{ID: 195}, PP: 182.23}, {ID: 196, Beatmap: Beatmap{ID: 196}, PP: 73.64},
				{ID: 197, Beatmap: Beatmap{ID: 197}, PP: 78.76}, {ID: 198, Beatmap: Beatmap{ID: 198}, PP: 765.47}, {ID: 199, Beatmap: Beatmap{ID: 199}, PP: 18.38}, {ID: 200, Beatmap: Beatmap{ID: 200}, PP: 417.73},
			},
			ExpectedCount: 100,
		},
		{
			Name: "Different scores on the same map",
			Scores: Scores{
				{ID: 1, Beatmap: Beatmap{ID: 1}, PP: 230.98}, {ID: 2, Beatmap: Beatmap{ID: 1}, PP: 307.03}, {ID: 3, Beatmap: Beatmap{ID: 1}, PP: 400.48},
				{ID: 4, Beatmap: Beatmap{ID: 1}, PP: 21.46}, {ID: 5, Beatmap: Beatmap{ID: 5}, PP: 98.31}, {ID: 6, Beatmap: Beatmap{ID: 1}, PP: 51.08},
			},
			ExpectedCount: 2,
		},
		{
			Name: "Trying to add the exact same score",
			Scores: Scores{
				{ID: 1, Beatmap: Beatmap{ID: 1}, PP: 230.98}, {ID: 1, Beatmap: Beatmap{ID: 1}, PP: 230.98}, {ID: 1, Beatmap: Beatmap{ID: 1}, PP: 230.98},
				{ID: 1, Beatmap: Beatmap{ID: 1}, PP: 230.98}, {ID: 1, Beatmap: Beatmap{ID: 1}, PP: 230.98}, {ID: 2, Beatmap: Beatmap{ID: 2}, PP: 230.98},
				{ID: 2, Beatmap: Beatmap{ID: 2}, PP: 240.98}, {ID: 3, Beatmap: Beatmap{ID: 3}, PP: 210.98}, {ID: 1, Beatmap: Beatmap{ID: 1}, PP: 230.98},
				{ID: 4, Beatmap: Beatmap{ID: 1}, PP: 280.98}, {ID: 2, Beatmap: Beatmap{ID: 2}, PP: 240.98}, {ID: 2, Beatmap: Beatmap{ID: 2}, PP: 240.98},
			},
			ExpectedCount: 3,
		},
	}
}

func TestRanking_AddScoreOrdered(t *testing.T) {
	for _, testCase := range getTestCases() {
		t.Run(testCase.Name, func(t *testing.T) {
			var rank Ranking
			for _, score := range testCase.Scores {
				rank.AddScore(score)
			}
			assertEqual(t, testCase.ExpectedCount, int(rank.count))
			for i := 1; i < int(rank.count); i++ {
				if rank.scores[i].PP > rank.scores[i-1].PP {
					t.Fatal("Rank is unordered")
				}
			}
		})
	}
}

func TestRanking_AddScoreUniqueScoreID(t *testing.T) {
	for _, testCase := range getTestCases() {
		t.Run(testCase.Name, func(t *testing.T) {
			var rank Ranking
			for _, score := range testCase.Scores {
				rank.AddScore(score)
			}
			var ids = make(map[int]bool)
			for _, score := range rank.scores {
				if score.ID == 0 {
					continue
				}
				if ids[score.ID] {
					t.Fatal("Non-Unique ID found")
				}
				ids[score.ID] = true
			}
			assertEqual(t, testCase.ExpectedCount, int(rank.count))
		})
	}
}

func TestRanking_AddScoreUniqueBeatmapID(t *testing.T) {
	for _, testCase := range getTestCases() {
		t.Run(testCase.Name, func(t *testing.T) {
			var rank Ranking
			for _, score := range testCase.Scores {
				rank.AddScore(score)
			}
			var ids = make(map[int]bool)
			for _, score := range rank.scores {
				if score.ID == 0 {
					continue
				}
				if ids[score.Beatmap.ID] {
					t.Fatal("Non-Unique ID found")
				}
				ids[score.Beatmap.ID] = true
			}
			assertEqual(t, testCase.ExpectedCount, int(rank.count))
		})
	}
}

func TestRanking_AddScore(t *testing.T) {
	var (
		scores = Scores{
			{ID: 1, Beatmap: Beatmap{ID: 1}, PP: 336.242},
			{ID: 2, Beatmap: Beatmap{ID: 2}, PP: 332.834},
			{ID: 3, Beatmap: Beatmap{ID: 3}, PP: 330.403},
			{ID: 4, Beatmap: Beatmap{ID: 4}, PP: 328.735},
			{ID: 5, Beatmap: Beatmap{ID: 5}, PP: 328.239},
		}

		arrExpectedTotalPP = []float64{
			336,
			316,
			298,
			282,
			267,
		}

		expectedTotalPP float64
	)
	log.Println(scores)

	var rank Ranking
	log.Println(&rank)

	for i, score := range scores {
		rank.AddScore(score)
		assertEqual(t, i+1, int(rank.count))

		for j := 1; j < int(rank.count); j++ {
			if rank.scores[j].PP > rank.scores[j-1].PP {
				t.Fatalf("Rank is unordered")
			}
		}

		expectedTotalPP += arrExpectedTotalPP[i]
		assertEqual(t, expectedTotalPP, rank.GetTotalPP())
	}
	log.Println(&rank)
}

func TestRanking_GetTotalPPSingle(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		const ExpectedTotalPP = float64(100)
		var rank = Ranking{
			count: 1,
			scores: [RankSize]Score{
				{PP: ExpectedTotalPP},
			},
		}

		assertEqual(t, ExpectedTotalPP, rank.GetTotalPP())
	})

	t.Run("Multiple", func(t *testing.T) {
		const ExpectedTotalPP = float64(336 + 316 + 298 + 282 + 267)
		var rank = Ranking{
			count: 5,
			scores: [RankSize]Score{
				{PP: 336.242},
				{PP: 332.834},
				{PP: 330.403},
				{PP: 328.735},
				{PP: 328.239},
			},
		}

		assertEqual(t, ExpectedTotalPP, rank.GetTotalPP())
	})
}

func TestRanking_Clone(t *testing.T) {
	var expected = Ranking{
		count: 5,
		scores: [RankSize]Score{
			{PP: 336.242},
			{PP: 332.834},
			{PP: 330.403},
			{PP: 328.735},
			{PP: 328.239},
		},
	}
	var actual = expected
	assertEqual(t, expected.count, actual.count)
	assert(t, &expected.scores != &actual.scores, "Pointing to the same ranking")

	actual.scores[0] = Score{}
	assert(t, expected.scores[0].PP > actual.scores[0].PP, "Messed up the ranking")

	for i := int8(1); i < actual.count; i++ {
		assertEqual(t, expected.scores[i].PP, actual.scores[i].PP)
	}
}

func BenchmarkRanking_AddScore(b *testing.B) {
	var testCases = getTestCases()
	for _, testCase := range testCases {
		b.Run(testCase.Name, func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				var rank Ranking
				for _, score := range testCase.Scores {
					rank.AddScore(score)
				}
			}
		})
	}
}

func BenchmarkRanking_GetTotalPP(b *testing.B) {
	// Build ranks from test cases
	var testCases = getTestCases()
	var ranks = make([]Ranking, 0, len(testCases))
	for _, testCase := range testCases {
		var rank Ranking
		for _, score := range testCase.Scores {
			rank.AddScore(score)
		}
		ranks = append(ranks, rank)
	}

	// Benchmark rankings
	for i, rank := range ranks {
		b.Run(testCases[i].Name, func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				_ = rank.GetTotalPP()
			}
		})
	}
}
