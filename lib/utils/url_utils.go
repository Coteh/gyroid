package utils

// IsURL determines whether supplied string is a URL
func IsURL(urlStr string) bool {
	return urlStr[0:4] == "http" || urlStr[0:5] == "https"
}
