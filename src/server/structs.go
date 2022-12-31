package server

type (
	// Expected before anything else is sent
	ServerLogin struct {
		Name string
	}

	// ID = 1
	ReqTrackUser struct {
		Username string
		Mode     string
	}
	ResTrackUser struct {
		Username string
		Message  string
		TotalPP  float64
	}
)
