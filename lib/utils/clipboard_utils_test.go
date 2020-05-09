package utils_test

import (
	"errors"
	"testing"

	"github.com/Coteh/gyroid/lib/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const CLIPBOARD_STRING_FIXTURE = "mockString"
const CLIPBOARD_HTTP_URL_STRING_FIXTURE = "http://www.google.com"
const CLIPBOARD_HTTPS_URL_STRING_FIXTURE = "https://www.google.com"
const CLIPBOARD_ERROR_FIXTURE = "clipboard error"

type MockClipboardManager struct {
	mock.Mock
}

func (m *MockClipboardManager) GetFromClipboard() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockClipboardManager) GetMostRecentlyAddedURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockClipboardManager) SetMostRecentlyAddedURL(url string) {
	m.Called()
}

func (m *MockClipboardManager) CopyToClipboard(text string) error {
	args := m.Called()
	return args.Error(0)
}

func TestIsURLInClipboardReturnsTrueForValidHTTPURL(t *testing.T) {
	clipboardManager := &MockClipboardManager{}
	clipboardManager.On("GetFromClipboard").Return(CLIPBOARD_HTTP_URL_STRING_FIXTURE, nil)
	assert.True(t, utils.IsURLInClipboard(clipboardManager))
}

func TestIsURLInClipboardReturnsTrueForValidHTTPSURL(t *testing.T) {
	clipboardManager := &MockClipboardManager{}
	clipboardManager.On("GetFromClipboard").Return(CLIPBOARD_HTTPS_URL_STRING_FIXTURE, nil)
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

func TestCheckURLMostRecentlyAddedReturnsTrueIfClipboardTextIsEqual(t *testing.T) {
	clipboardManager := &MockClipboardManager{}
	expectedURL := CLIPBOARD_HTTP_URL_STRING_FIXTURE
	clipboardManager.On("GetFromClipboard").Return(expectedURL, nil)
	clipboardManager.On("GetMostRecentlyAddedURL").Return(expectedURL)
	assert.True(t, utils.CheckURLMostRecentlyAdded(clipboardManager))
}

func TestCheckURLMostRecentlyAddedReturnsFalseIfClipboardTextIsNotEqual(t *testing.T) {
	clipboardManager := &MockClipboardManager{}
	clipboardManager.On("GetFromClipboard").Return(CLIPBOARD_STRING_FIXTURE, nil)
	clipboardManager.On("GetMostRecentlyAddedURL").Return(CLIPBOARD_HTTP_URL_STRING_FIXTURE)
	assert.False(t, utils.CheckURLMostRecentlyAdded(clipboardManager))
}

func TestCheckURLMostRecentlyAddedReturnsFalseIfMostRecentlyAddedURLIsEmpty(t *testing.T) {
	clipboardManager := &MockClipboardManager{}
	clipboardManager.On("GetFromClipboard").Return(CLIPBOARD_HTTP_URL_STRING_FIXTURE, nil)
	clipboardManager.On("GetMostRecentlyAddedURL").Return("")
	assert.False(t, utils.CheckURLMostRecentlyAdded(clipboardManager))
}

func TestCheckURLMostRecentlyAddedReturnsFalseIfErrorGettingClipboardText(t *testing.T) {
	clipboardManager := &MockClipboardManager{}
	expectedURL := CLIPBOARD_HTTP_URL_STRING_FIXTURE
	clipboardManager.On("GetFromClipboard").Return("", errors.New(CLIPBOARD_ERROR_FIXTURE))
	clipboardManager.On("GetMostRecentlyAddedURL").Return(expectedURL)
	assert.False(t, utils.CheckURLMostRecentlyAdded(clipboardManager))
}
