package utils_test

import (
	"fmt"
	"github.com/Coteh/gyroid/lib/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenAuthURLOpensAuthURL(t *testing.T) {
	requestTokenFixture := "test"
	redirectURIFixture := "my-url"

	expectedURL := fmt.Sprintf("https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s", requestTokenFixture, redirectURIFixture)

	testOpenBrowser := func(url string) {
		assert.Equal(t, expectedURL, url)
	}

	utils.OpenAuthURL(requestTokenFixture, redirectURIFixture, testOpenBrowser)
}
