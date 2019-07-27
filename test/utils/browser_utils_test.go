package utils_test

import (
	"fmt"
	"github.com/Coteh/gyroid/lib/utils"
	"testing"
)

func TestOpenAuthURLOpensAuthURL(t *testing.T) {
	requestTokenFixture := "test"
	redirectURIFixture := "my-url"

	expectedURL := fmt.Sprintf("https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s", requestTokenFixture, redirectURIFixture)

	testOpenBrowser := func(url string) {
		if url != expectedURL {
			t.FailNow()
		}
	}

	utils.OpenAuthURL(requestTokenFixture, redirectURIFixture, testOpenBrowser)
}
