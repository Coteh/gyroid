package utils

import (
	"net/url"
)

// IsURL determines whether supplied string is a URL
// https://stackoverflow.com/a/55551215/9292680
func IsURL(urlStr string) bool {
	url, err := url.Parse(urlStr)
	return err == nil && url.Host != "" && url.Scheme != ""
}
