package actions_test

import (
	"errors"
	"testing"

	"github.com/Coteh/gyroid/lib/actions"
	"github.com/Coteh/gyroid/lib/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteArticleCallsModifyWithCorrectParams(t *testing.T) {
	mockClient := &PocketClientSimpleActionMock{}
	expectedAction := models.PocketAction{
		Action: "delete",
		ItemID: ARTICLE_ID_FIXTURE,
	}
	expectedActionArr := make([]models.PocketAction, 1)
	expectedActionArr[0] = expectedAction
	expectedParams := models.PocketModify{
		Actions: expectedActionArr,
	}
	mockClient.On("Modify", expectedParams)
	actions.DeleteArticle(mockClient, ARTICLE_ID_FIXTURE)
}

func TestDeleteArticleReturnsTrueOnSuccess(t *testing.T) {
	mockClient := &PocketClientSimpleActionMock{}
	mockClient.On("Modify", mock.Anything)
	result, _ := actions.DeleteArticle(mockClient, ARTICLE_ID_FIXTURE)
	assert.True(t, result)
}

func TestDeleteArticleReturnsFalseOnClientFailure(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("Modify", mock.Anything).Return(nil, errors.New(MOCK_ERROR_STRING))
	result, _ := actions.DeleteArticle(mockClient, ARTICLE_ID_FIXTURE)
	assert.False(t, result)
}

func TestDeleteArticleReturnsClientErrorOnClientFailure(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("Modify", mock.Anything).Return(nil, errors.New(MOCK_ERROR_STRING))
	_, err := actions.DeleteArticle(mockClient, ARTICLE_ID_FIXTURE)
	assert.Equal(t, MOCK_ERROR_STRING, err.Error())
}
