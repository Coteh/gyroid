package utils_test

import (
	"testing"

	"github.com/Coteh/gyroid/lib/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsURLReturnsTrueForValidHttpUrl(t *testing.T) {
	url := "http://www.google.com"
	assert.True(t, utils.IsURL(url))
}

func TestIsURLReturnsTrueForValidHttpsUrl(t *testing.T) {
	url := "https://www.google.com"
	assert.True(t, utils.IsURL(url))
}

func TestIsURLReturnsFalseForInvalidUrl(t *testing.T) {
	url := "invalid-url"
	assert.False(t, utils.IsURL(url))
}
