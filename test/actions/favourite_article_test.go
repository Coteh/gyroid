package actions_test

import (
	"github.com/Coteh/gyroid/lib/actions"
	"github.com/Coteh/gyroid/lib/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFavouriteArticleCallsModifyWithCorrectParams(t *testing.T) {
	mockClient := &PocketClientSimpleActionMock{}
	expectedAction := models.PocketAction{
		Action: "favorite",
		ItemID: ARTICLE_ID_FIXTURE,
	}
	expectedActionArr := make([]models.PocketAction, 1)
	expectedActionArr[0] = expectedAction
	expectedParams := models.PocketModify{
		Actions: expectedActionArr,
	}
	mockClient.On("Modify", expectedParams)
	actions.FavouriteArticle(mockClient, ARTICLE_ID_FIXTURE)
}

func TestFavouriteArticleReturnsTrueOnSuccess(t *testing.T) {
	mockClient := &PocketClientSimpleActionMock{}
	mockClient.On("Modify", mock.Anything)
	result, _ := actions.FavouriteArticle(mockClient, ARTICLE_ID_FIXTURE)
	assert.True(t, result)
}

func TestFavouriteArticleReturnsFalseOnClientFailure(t *testing.T) {
	mockClient := &FailingPocketClientMock{}
	result, _ := actions.FavouriteArticle(mockClient, ARTICLE_ID_FIXTURE)
	assert.False(t, result)
}

func TestFavouriteArticleReturnsClientErrorOnClientFailure(t *testing.T) {
	mockClient := &FailingPocketClientMock{}
	_, err := actions.FavouriteArticle(mockClient, ARTICLE_ID_FIXTURE)
	assert.Equal(t, MOCK_ERROR_STRING, err.Error())
}
