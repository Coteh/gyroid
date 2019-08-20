package actions_test

import (
	"testing"

	"github.com/Coteh/gyroid/lib/actions"
	"github.com/Coteh/gyroid/lib/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const URL_FIXTURE = "http://www.google.com"

func TestAddArticleCallsAddWithCorrectParams(t *testing.T) {
	expectedParams := models.PocketAdd{
		Url: URL_FIXTURE,
	}
	mockClient := CreateSuccessfulPocketClientMock("Add", CreateSuccessfulAddResult(expectedParams.Url), expectedParams)
	actions.AddArticle(mockClient, URL_FIXTURE)
}

func TestAddArticleReturnsItemOnSuccess(t *testing.T) {
	expectedItem := new(models.AddedArticleResult)
	expectedItem.NormalURL = URL_FIXTURE
	mockClient := CreateSuccessfulPocketClientMock("Add", CreateSuccessfulAddResult(URL_FIXTURE), mock.Anything)
	result, _ := actions.AddArticle(mockClient, URL_FIXTURE)
	assert.Equal(t, expectedItem, result)
}

func TestAddArticleReturnsErrorOnClientFailure(t *testing.T) {
	mockClient := CreateFailingPocketClientMock("Add", mock.Anything)
	result, err := actions.AddArticle(mockClient, URL_FIXTURE)
	assert.Nil(t, result)
	assert.Equal(t, MOCK_ERROR_STRING, err.Error())
}
