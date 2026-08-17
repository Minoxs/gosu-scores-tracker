package osu

import (
	"net/http"
	"time"
)

const (
	BaseURL  = "https://osu.ppy.sh"
	OAuthURL = BaseURL + "/oauth"
	ApiV2    = BaseURL + "/api/v2"

	GET  = "GET"
	AUTH = "Authorization"
	JSON = "application/json"
)

// apiClient is the client used for requests in this package. Its transport paces
// every request through globalPacer to respect the osu! API rate limit.
var apiClient = &http.Client{
	Timeout:   30 * time.Second,
	Transport: &throttledTransport{base: http.DefaultTransport, pacer: globalPacer},
}
