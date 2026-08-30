package tracker

import "github.com/minoxs/gosu-api/pkg/gosu"

// DefaultProvider is the simplest form of a provider
// which just returns the given token.
type DefaultProvider struct {
	Token *gosu.GuestToken
}

// GetToken returns the token in the provider
func (p *DefaultProvider) GetToken() *gosu.GuestToken {
	return p.Token
}
