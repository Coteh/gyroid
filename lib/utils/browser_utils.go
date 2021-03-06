package utils

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

// OpenAuthURL opens auth url
func OpenAuthURL(code string, redirectURI string, browserOpenFunc func(string)) {
	url := fmt.Sprintf("https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s", code, redirectURI)

	browserOpenFunc(url)
}
