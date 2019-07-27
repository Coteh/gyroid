package utils

import (
	"github.com/atotto/clipboard"
)

type ClipboardManager interface {
	GetFromClipboard() (string, error)
}

type ClipboardManagerImpl struct {
}

func (c *ClipboardManagerImpl) GetFromClipboard() (string, error) {
	return clipboard.ReadAll()
}

func IsURLInClipboard(clipboardManager ClipboardManager) bool {
	result, err := clipboardManager.GetFromClipboard()
	return err == nil && IsURL(result)
}

func GetFromClipboard(clipboardManager ClipboardManager) (string, error) {
	return clipboardManager.GetFromClipboard()
}
