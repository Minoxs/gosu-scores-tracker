package player

import "time"

// Profile is a user's public osu! profile as returned by the users endpoint.
type Profile struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	CountryCode string    `json:"country_code"`
	AvatarURL   string    `json:"avatar_url"`
	CoverURL    string    `json:"cover_url"`
	IsSupporter bool      `json:"is_supporter"`
	IsOnline    bool      `json:"is_online"`
	JoinDate    time.Time `json:"join_date"`

	Country struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"country"`

	Cover struct {
		URL string `json:"url"`
	} `json:"cover"`

	Statistics struct {
		PP          float64 `json:"pp"`
		HitAccuracy float64 `json:"hit_accuracy"`
		PlayCount   int     `json:"play_count"`
		GlobalRank  *int    `json:"global_rank"`
		CountryRank *int    `json:"country_rank"`
		Rank        struct {
			Country *int `json:"country"`
		} `json:"rank"`
	} `json:"statistics"`
}

// Cover image URL, preferring the nested cover object the API populates.
func (p *Profile) Banner() string {
	if p.Cover.URL != "" {
		return p.Cover.URL
	}
	return p.CoverURL
}

// GlobalRank is the all-time global rank, nil when the user is unranked.
func (p *Profile) GlobalRank() *int {
	return p.Statistics.GlobalRank
}

// CountryRank falls back to the legacy rank.country field when the flat one is absent.
func (p *Profile) CountryRank() *int {
	if p.Statistics.CountryRank != nil {
		return p.Statistics.CountryRank
	}
	return p.Statistics.Rank.Country
}

// Country2 is the ISO alpha-2 code, preferring the nested country object.
func (p *Profile) Country2() string {
	if p.Country.Code != "" {
		return p.Country.Code
	}
	return p.CountryCode
}
