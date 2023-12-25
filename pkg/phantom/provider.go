package phantom

import "osu-phantom/pkg/osu"

type DefaultProvider struct {
	Token *osu.GuestToken
}

func (p *DefaultProvider) GetToken() *osu.GuestToken {
	return p.Token
}
