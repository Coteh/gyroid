package actions_test

import (
	"github.com/Coteh/gyroid/lib/actions"
	"github.com/Coteh/gyroid/lib/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const URL_FIXTURE = "http://www.google.com"

type PocketArticleAddArticleMock struct {
	PocketClientMock
}

func (m *PocketArticleAddArticleMock) Add(params models.PocketAdd) (*models.PocketAddResult, error) {
	m.Called(params)
	mockItem := make(map[string]interface{})
	mockItem["normal_url"] = params.Url
	mockResult := &models.PocketAddResult{
		Status: 0,
		Item:   mockItem,
	}
	return mockResult, nil
}

func TestAddArticleCallsAddWithCorrectParams(t *testing.T) {
	mockClient := &PocketArticleAddArticleMock{}
	expectedParams := models.PocketAdd{
		Url: URL_FIXTURE,
	}
	mockClient.On("Add", expectedParams)
	actions.AddArticle(mockClient, URL_FIXTURE)
}

func TestAddArticleReturnsItemOnSuccess(t *testing.T) {
	mockClient := &PocketArticleAddArticleMock{}
	expectedItem := make(map[string]interface{})
	expectedItem["normal_url"] = URL_FIXTURE
	mockClient.On("Add", mock.Anything)
	result, _ := actions.AddArticle(mockClient, URL_FIXTURE)
	assert.Equal(t, expectedItem, result)
}

func TestAddArticleReturnsNilOnClientFailure(t *testing.T) {
	mockClient := &FailingPocketClientMock{}
	result, _ := actions.AddArticle(mockClient, URL_FIXTURE)
	assert.Nil(t, result)
}

func TestAddArticleReturnsErrorOnClientFailure(t *testing.T) {
	mockClient := &FailingPocketClientMock{}
	_, err := actions.AddArticle(mockClient, URL_FIXTURE)
	assert.Equal(t, MOCK_ERROR_STRING, err.Error())
}
