package osu

import (
	"github.com/minoxs/osu-phantom/pkg/osu/optimization"
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

var (
	// apiClient is the client used for requests in this package
	apiClient = &http.Client{
		Timeout: 30 * time.Second,
	}

	// cache is used to store beatmaps without having to ask the API every time. Disabled by default.
	// Call SetCacheOptions to enable it.
	// TODO REMOVE FROM PACKAGE SCOPE
	cache = &optimization.BeatmapCache{
		MaxUnitSize: 0,
		CacheSize:   0,
	}
)

func init() {
	cache.Init()
}

func SetCacheOptions(MaxUnitSize uint32, CacheSize uint32) {
	cache.MaxUnitSize = MaxUnitSize
	cache.CacheSize = CacheSize
}
