package utils

import (
	"github.com/atotto/clipboard"
)

// ClipboardManager is an interface for actions with the clipboard
type ClipboardManager interface {
	GetFromClipboard() (string, error)
	GetMostRecentlyAddedURL() string
	SetMostRecentlyAddedURL(url string)
}

// ClipboardManagerImpl is an implementation of ClipboardManager that uses atotto's clipboard library
type ClipboardManagerImpl struct {
	AddedURL string
}

// GetFromClipboard gets the clipboard text using atotto's clipboard library
func (c *ClipboardManagerImpl) GetFromClipboard() (string, error) {
	return clipboard.ReadAll()
}

// GetMostRecentlyAddedURL gets the URL that was most recently added
func (c *ClipboardManagerImpl) GetMostRecentlyAddedURL() string {
	return c.AddedURL
}

// SetMostRecentlyAddedURL sets the URL that was most recently added
func (c *ClipboardManagerImpl) SetMostRecentlyAddedURL(url string) {
	c.AddedURL = url
}

// IsURLInClipboard checks to see if the clipboard text is actually a URL
func IsURLInClipboard(clipboardManager ClipboardManager) bool {
	result, err := clipboardManager.GetFromClipboard()
	return err == nil && IsURL(result)
}

// CheckURLMostRecentlyAdded checks to see if the text in clipboard is the same as the URL that was just added
func CheckURLMostRecentlyAdded(clipboardManager ClipboardManager) bool {
	result, err := clipboardManager.GetFromClipboard()
	return err == nil && result == clipboardManager.GetMostRecentlyAddedURL()
}

// GetFromClipboard gets clipboard text from the clipboard manager
func GetFromClipboard(clipboardManager ClipboardManager) (string, error) {
	return clipboardManager.GetFromClipboard()
}
