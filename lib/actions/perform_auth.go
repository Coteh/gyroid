package actions

import (
	pocketConnector "github.com/Coteh/gyroid/lib/connector"
	"github.com/Coteh/gyroid/lib/utils"
	"time"
)

// PerformAuth performs authentication with Pocket to retrieve an access token for the user
func PerformAuth(client pocketConnector.PocketConnector, delayMs time.Duration, redirectURI string, browserOpenFunc func(string)) (string, error) {
	code, err := client.RequestOAuthCode(redirectURI)
	if err != nil {
		return "", err
	}

	utils.OpenAuthURL(code, redirectURI, browserOpenFunc)

	time.Sleep(delayMs * time.Millisecond)

	accessToken, err := client.Authorize(code)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
