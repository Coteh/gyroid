package utils

import (
	"os"
	"runtime"
)

// UserHomeDir gets the user's home directory based on OS
// https://stackoverflow.com/a/7922977/9292680
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
