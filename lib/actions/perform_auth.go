package actions

import (
	"time"

	"github.com/Coteh/gyroid/lib/connector"
	"github.com/Coteh/gyroid/lib/utils"

	"net/http"
	"net/http/httptest"
)

// PerformAuth performs authentication with Pocket to retrieve an access token for the user
func PerformAuth(client connector.PocketConnector, delayMs time.Duration, browserOpenFunc func(string)) (string, error) {
	ts := launchRedirectServer()
	defer ts.Close()

	code, err := client.RequestOAuthCode(ts.URL)
	if err != nil {
		return "", err
	}

	utils.OpenAuthURL(code, ts.URL, browserOpenFunc)

	time.Sleep(delayMs * time.Millisecond)

	accessToken, err := client.Authorize(code)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func launchRedirectServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("Pocket user has been authorized. Close this window and go back to the app."))
		}),
	)
}
