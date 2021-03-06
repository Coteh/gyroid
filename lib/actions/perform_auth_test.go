package actions_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/Coteh/gyroid/lib/actions"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const CODE_FIXTURE = "oauthCode"
const REDIRECT_URI_FIXTURE = "http://localhost:8080"
const ACCESS_TOKEN_FIXTURE = "accessToken"
const DELAY_MILLISECONDS = 0
const EXPECTED_AUTH_URL_PREFIX = "https://getpocket.com/auth/authorize?request_token=" + CODE_FIXTURE

func openURLStub(redirectUri string) {}

func TestPerformAuthPerformsAuth(t *testing.T) {
	mockClient := CreatePocketClientMock()
	mockClient.On("RequestOAuthCode", mock.Anything).Return(CODE_FIXTURE, nil)
	mockClient.On("Authorize", mock.Anything).Return(ACCESS_TOKEN_FIXTURE, nil)
	result, _ := actions.PerformAuth(mockClient, DELAY_MILLISECONDS, openURLStub)
	assert.Equal(t, ACCESS_TOKEN_FIXTURE, result)
}

func TestPerformAuthCallsAuthorizeWithCorrectParams(t *testing.T) {
	mockClient := CreatePocketClientMock()
	mockClient.On("RequestOAuthCode", mock.Anything).Return(CODE_FIXTURE, nil)
	mockClient.On("Authorize", CODE_FIXTURE).Return(ACCESS_TOKEN_FIXTURE, nil)
	actions.PerformAuth(mockClient, DELAY_MILLISECONDS, openURLStub)
}

func TestPerformAuthReturnsErrorOnAuthorizeFailure(t *testing.T) {
	mockClient := CreatePocketClientMock()
	mockClient.On("RequestOAuthCode", mock.Anything).Return(CODE_FIXTURE, nil)
	mockClient.On("Authorize", CODE_FIXTURE).Return("", errors.New(MOCK_ERROR_STRING))
	result, err := actions.PerformAuth(mockClient, DELAY_MILLISECONDS, openURLStub)
	assert.Empty(t, result)
	assert.Equal(t, MOCK_ERROR_STRING, err.Error())
}

func TestPerformAuthReturnsErrorOnRequestOAuthCodeFailure(t *testing.T) {
	mockClient := CreatePocketClientMock()
	mockClient.On("RequestOAuthCode", mock.Anything).Return("", errors.New(MOCK_ERROR_STRING))
	result, err := actions.PerformAuth(mockClient, DELAY_MILLISECONDS, openURLStub)
	assert.Empty(t, result)
	assert.Equal(t, MOCK_ERROR_STRING, err.Error())
}

func TestPerformAuthCallsOpenURLCorrectly(t *testing.T) {
	openURLMock := func(url string) {
		assert.True(t, strings.HasPrefix(url, EXPECTED_AUTH_URL_PREFIX))
	}
	mockClient := CreatePocketClientMock()
	mockClient.On("RequestOAuthCode", mock.Anything).Return(CODE_FIXTURE, nil)
	mockClient.On("Authorize", mock.Anything).Return(ACCESS_TOKEN_FIXTURE, nil)
	actions.PerformAuth(mockClient, DELAY_MILLISECONDS, openURLMock)
}
