package utils

import (
	"net/url"
)

func URLValidator(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || !parsedURL.IsAbs() {
		return false
	}
	return true
}
