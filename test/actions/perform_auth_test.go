package actions_test

import (
	"errors"
	"testing"

	"github.com/Coteh/gyroid/lib/actions"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const CODE_FIXTURE = "oauthCode"
const REDIRECT_URI_FIXTURE = "http://localhost:8080"
const ACCESS_TOKEN_FIXTURE = "accessToken"
const DELAY_MILLISECONDS = 0
const EXPECTED_AUTH_URL = "https://getpocket.com/auth/authorize?request_token=" + CODE_FIXTURE + "&redirect_uri=" + REDIRECT_URI_FIXTURE

func openURLStub(redirectUri string) {}

func TestPerformAuthPerformsAuth(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("RequestOAuthCode", mock.Anything).Return(CODE_FIXTURE, nil)
	mockClient.On("Authorize", mock.Anything).Return(ACCESS_TOKEN_FIXTURE, nil)
	result, _ := actions.PerformAuth(mockClient, DELAY_MILLISECONDS, REDIRECT_URI_FIXTURE, openURLStub)
	assert.Equal(t, ACCESS_TOKEN_FIXTURE, result)
}

func TestPerformAuthCallsAuthorizeWithCorrectParams(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("RequestOAuthCode", REDIRECT_URI_FIXTURE).Return(CODE_FIXTURE, nil)
	mockClient.On("Authorize", CODE_FIXTURE).Return(ACCESS_TOKEN_FIXTURE, nil)
	actions.PerformAuth(mockClient, DELAY_MILLISECONDS, REDIRECT_URI_FIXTURE, openURLStub)
}

func TestPerformAuthReturnsEmptyStringOnAuthorizeFailure(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("RequestOAuthCode", REDIRECT_URI_FIXTURE).Return(CODE_FIXTURE, nil)
	mockClient.On("Authorize", CODE_FIXTURE).Return("", errors.New(MOCK_ERROR_STRING))
	result, _ := actions.PerformAuth(mockClient, DELAY_MILLISECONDS, REDIRECT_URI_FIXTURE, openURLStub)
	assert.Empty(t, result)
}

func TestPerformAuthReturnsErrorOnAuthorizeFailure(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("RequestOAuthCode", REDIRECT_URI_FIXTURE).Return(CODE_FIXTURE, nil)
	mockClient.On("Authorize", CODE_FIXTURE).Return("", errors.New(MOCK_ERROR_STRING))
	_, err := actions.PerformAuth(mockClient, DELAY_MILLISECONDS, REDIRECT_URI_FIXTURE, openURLStub)
	assert.Equal(t, MOCK_ERROR_STRING, err.Error())
}

func TestPerformAuthReturnsEmptyStringOnOAuthFailure(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("RequestOAuthCode", REDIRECT_URI_FIXTURE).Return("", errors.New(MOCK_ERROR_STRING))
	result, _ := actions.PerformAuth(mockClient, DELAY_MILLISECONDS, REDIRECT_URI_FIXTURE, openURLStub)
	assert.Empty(t, result)
}

func TestPerformAuthReturnsErrorOnOAuthFailure(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("RequestOAuthCode", REDIRECT_URI_FIXTURE).Return("", errors.New(MOCK_ERROR_STRING))
	_, err := actions.PerformAuth(mockClient, DELAY_MILLISECONDS, REDIRECT_URI_FIXTURE, openURLStub)
	assert.Equal(t, MOCK_ERROR_STRING, err.Error())
}

func TestPerformAuthCallsOpenURLCorrectly(t *testing.T) {
	openURLMock := func(redirectUri string) {
		assert.Equal(t, EXPECTED_AUTH_URL, redirectUri)
	}
	mockClient := &PocketClientMock{}
	mockClient.On("RequestOAuthCode", mock.Anything).Return(CODE_FIXTURE, nil)
	mockClient.On("Authorize", mock.Anything).Return(ACCESS_TOKEN_FIXTURE, nil)
	actions.PerformAuth(mockClient, DELAY_MILLISECONDS, REDIRECT_URI_FIXTURE, openURLMock)
}
