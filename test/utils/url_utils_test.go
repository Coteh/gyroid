package utils_test

import (
	"github.com/Coteh/gyroid/lib/utils"
	"testing"
)

func TestIsURLReturnsTrueForValidHttpUrl(t *testing.T) {
	url := "http://www.google.com"

	if !utils.IsURL(url) {
		t.FailNow()
	}
}

func TestIsURLReturnsTrueForValidHttpsUrl(t *testing.T) {
	url := "https://www.google.com"

	if !utils.IsURL(url) {
		t.FailNow()
	}
}

func TestIsURLReturnsFalseForInvalidUrl(t *testing.T) {
	url := "invalid-url"

	if utils.IsURL(url) {
		t.FailNow()
	}
}
