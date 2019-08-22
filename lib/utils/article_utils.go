package utils

const wpm = 265

// CalculateExpectedReadTime calculates approx. read time in minutes based on Medium's estimate
// https://help.medium.com/hc/en-us/articles/214991667-Read-time
func CalculateExpectedReadTime(wordCount int) int {
	return wordCount / wpm
}
