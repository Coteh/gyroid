package config_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/Coteh/gyroid/lib/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const ERROR_STRING_FIXTURE = "file i/o error"

type FailingReader struct {
}

func (r *FailingReader) Read(b []byte) (int, error) {
	return 0, errors.New(ERROR_STRING_FIXTURE)
}

type MockWriter struct {
	mock.Mock
}

func (m *MockWriter) Write(p []byte) (int, error) {
	args := m.Called(p)

	return args.Int(0), args.Error(1)
}

func TestReadConfigReturnsValidConfigIfParsed(t *testing.T) {
	reader := strings.NewReader("clipboard: true")

	configResult, err := config.ReadConfig(reader)
	assert.Nil(t, err)

	assert.Equal(t, &config.Config{Clipboard: true}, configResult)
}

func TestReadConfigReturnsErrorIfParsingError(t *testing.T) {
	reader := &FailingReader{}

	configResult, err := config.ReadConfig(reader)

	assert.Nil(t, configResult)
	assert.Equal(t, ERROR_STRING_FIXTURE, err.Error())
}

func TestReadConfigReturnsDefaultsIfNoConfigData(t *testing.T) {
	reader := strings.NewReader("")

	configResult, _ := config.ReadConfig(reader)

	assert.Equal(t, &config.Config{}, configResult)
}

func TestWriteConfigCanCallWriteIfValidConfigPassed(t *testing.T) {
	writer := &MockWriter{}

	configObj := &config.Config{Clipboard: true}

	writer.On("Write", mock.Anything).Return(10, nil)
	config.WriteConfig(configObj, writer)
}

func TestWriteConfigDoesNotCallWriteIfNilConfigPassed(t *testing.T) {
	writer := &MockWriter{}

	config.WriteConfig(nil, writer)

	writer.AssertNotCalled(t, "Write", mock.Anything)
}

func TestWriteConfigReturnsErrorIfNilConfigPassed(t *testing.T) {
	writer := &MockWriter{}

	err := config.WriteConfig(nil, writer)

	assert.Equal(t, config.NO_CONFIG_WRITE_ERROR, err.Error())
}

func TestWriteConfigReturnsErrorIfWritingError(t *testing.T) {
	writer := &MockWriter{}

	configObj := &config.Config{}

	writer.On("Write", mock.Anything).Return(0, errors.New(ERROR_STRING_FIXTURE))
	err := config.WriteConfig(configObj, writer)

	assert.Equal(t, ERROR_STRING_FIXTURE, err.Error())
}

func TestWriteConfigReturnsErrorIfNoBytesWritten(t *testing.T) {
	writer := &MockWriter{}

	configObj := &config.Config{}

	writer.On("Write", mock.Anything).Return(0, nil)
	err := config.WriteConfig(configObj, writer)

	assert.Equal(t, config.NO_BYTES_WRITTEN_ERROR, err.Error())
}
