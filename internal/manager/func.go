package manager

import (
	"net/url"
)

func ValidateURL(s string) (url.URL, bool) {
	parsed, err := url.Parse(s)
	return *parsed, err == nil && parsed.Scheme != "" && parsed.Host != ""
}
