package phantom

import "github.com/minoxs/osu-phantom/pkg/osu"

// DefaultProvider is the simplest form of a provider
// which just returns the given token.
type DefaultProvider struct {
	Token *osu.GuestToken
}

// GetToken returns the token in the provider
func (p *DefaultProvider) GetToken() *osu.GuestToken {
	return p.Token
}
