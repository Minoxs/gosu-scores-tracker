package crosu

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
