package actions_test

import (
	"testing"

	"github.com/Coteh/gyroid/lib/actions"
	"github.com/Coteh/gyroid/lib/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUnfavouriteArticleCallsModifyWithCorrectParams(t *testing.T) {
	expectedAction := models.PocketAction{
		Action: "unfavorite",
		ItemID: ARTICLE_ID_FIXTURE,
	}
	expectedActionArr := make([]models.PocketAction, 1)
	expectedActionArr[0] = expectedAction
	expectedParams := models.PocketModify{
		Actions: expectedActionArr,
	}
	mockClient := CreateSuccessfulPocketClientMock("Modify", CreateSuccessfulModifyResult(), expectedParams)
	actions.UnfavouriteArticle(mockClient, ARTICLE_ID_FIXTURE)
}

func TestUnfavouriteArticleReturnsTrueOnSuccess(t *testing.T) {
	mockClient := CreateSuccessfulPocketClientMock("Modify", CreateSuccessfulModifyResult(), mock.Anything)
	result, _ := actions.UnfavouriteArticle(mockClient, ARTICLE_ID_FIXTURE)
	assert.True(t, result)
}

func TestUnfavouriteArticleReturnsFalseOnClientFailure(t *testing.T) {
	mockClient := CreateFailingPocketClientMock("Modify", mock.Anything)
	result, _ := actions.UnfavouriteArticle(mockClient, ARTICLE_ID_FIXTURE)
	assert.False(t, result)
}

func TestUnfavouriteArticleReturnsClientErrorOnClientFailure(t *testing.T) {
	mockClient := CreateFailingPocketClientMock("Modify", mock.Anything)
	_, err := actions.UnfavouriteArticle(mockClient, ARTICLE_ID_FIXTURE)
	assert.Equal(t, MOCK_ERROR_STRING, err.Error())
}
