package url

import (
	"net/url"
)

// URL embeds net/url and adds UnmarshalText
type URL struct {
	url.URL
}

// UnmarshalText parses a URL
func (u *URL) UnmarshalText(text []byte) error {
	return u.UnmarshalBinary(text)
}
