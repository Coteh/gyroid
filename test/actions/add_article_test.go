package actions_test

import (
	"errors"
	"testing"

	"github.com/Coteh/gyroid/lib/actions"
	"github.com/Coteh/gyroid/lib/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const URL_FIXTURE = "http://www.google.com"

func TestAddArticleCallsAddWithCorrectParams(t *testing.T) {
	mockClient := &PocketClientMock{}
	expectedParams := models.PocketAdd{
		Url: URL_FIXTURE,
	}
	mockClient.On("Add", expectedParams).Return(CreateSuccessfulAddResult(expectedParams.Url), nil)
	actions.AddArticle(mockClient, URL_FIXTURE)
}

func TestAddArticleReturnsItemOnSuccess(t *testing.T) {
	mockClient := &PocketClientMock{}
	expectedItem := make(map[string]interface{})
	expectedItem["normal_url"] = URL_FIXTURE
	mockClient.On("Add", mock.Anything).Return(CreateSuccessfulAddResult(URL_FIXTURE), nil)
	result, _ := actions.AddArticle(mockClient, URL_FIXTURE)
	assert.Equal(t, expectedItem, result)
}

func TestAddArticleReturnsNilOnClientFailure(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("Add", mock.Anything).Return(nil, errors.New(MOCK_ERROR_STRING))
	result, _ := actions.AddArticle(mockClient, URL_FIXTURE)
	assert.Nil(t, result)
}

func TestAddArticleReturnsErrorOnClientFailure(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("Add", mock.Anything).Return(nil, errors.New(MOCK_ERROR_STRING))
	_, err := actions.AddArticle(mockClient, URL_FIXTURE)
	assert.Equal(t, MOCK_ERROR_STRING, err.Error())
}
