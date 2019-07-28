package utils_test

import (
	"errors"
	"testing"

	"github.com/Coteh/gyroid/lib/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const CLIPBOARD_STRING_FIXTURE = "mockString"
const CLIPBOARD_ERROR_FIXTURE = "clipboard error"

type MockClipboardManager struct {
	mock.Mock
	mockStr string
}

func (m *MockClipboardManager) GetFromClipboard() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestIsURLInClipboardReturnsTrueForValidHTTPURL(t *testing.T) {
	clipboardManager := &MockClipboardManager{}
	clipboardManager.On("GetFromClipboard").Return("http://www.google.com", nil)
	assert.True(t, utils.IsURLInClipboard(clipboardManager))
}

func TestIsURLInClipboardReturnsTrueForValidHTTPSURL(t *testing.T) {
	clipboardManager := &MockClipboardManager{}
	clipboardManager.On("GetFromClipboard").Return("https://www.google.com", nil)
	assert.True(t, utils.IsURLInClipboard(clipboardManager))
}

func TestIsURLInClipboardReturnsFalseForInvalidURL(t *testing.T) {
	clipboardManager := &MockClipboardManager{}
	clipboardManager.On("GetFromClipboard").Return("httpsinvalidurl", nil)
	assert.False(t, utils.IsURLInClipboard(clipboardManager))
}

func TestIsURLInClipboardReturnsFalseForClipboardError(t *testing.T) {
	clipboardManager := &MockClipboardManager{}
	clipboardManager.On("GetFromClipboard").Return("", errors.New(CLIPBOARD_ERROR_FIXTURE))
	assert.False(t, utils.IsURLInClipboard(clipboardManager))
}

func TestGetFromClipboardReturnsStringFromClipboardManager(t *testing.T) {
	clipboardManager := &MockClipboardManager{}
	clipboardManager.On("GetFromClipboard").Return(CLIPBOARD_STRING_FIXTURE, nil)
	result, _ := utils.GetFromClipboard(clipboardManager)
	assert.Equal(t, CLIPBOARD_STRING_FIXTURE, result)
}

func TestGetFromClipboardReturnsErrorIfClipboardManagerError(t *testing.T) {
	clipboardManager := &MockClipboardManager{}
	clipboardManager.On("GetFromClipboard").Return("", errors.New(CLIPBOARD_ERROR_FIXTURE))
	_, err := utils.GetFromClipboard(clipboardManager)
	assert.Equal(t, CLIPBOARD_ERROR_FIXTURE, err.Error())
}
