package actions_test

import (
	"errors"
	"testing"

	"github.com/Coteh/gyroid/lib/actions"
	"github.com/Coteh/gyroid/lib/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArchiveArticleCallsModifyWithCorrectParams(t *testing.T) {
	mockClient := &PocketClientMock{}
	expectedAction := models.PocketAction{
		Action: "archive",
		ItemID: ARTICLE_ID_FIXTURE,
	}
	expectedActionArr := make([]models.PocketAction, 1)
	expectedActionArr[0] = expectedAction
	expectedParams := models.PocketModify{
		Actions: expectedActionArr,
	}
	mockClient.On("Modify", expectedParams).Return(CreateSuccessfulModifyResult(), nil)
	actions.ArchiveArticle(mockClient, ARTICLE_ID_FIXTURE)
}

func TestArchiveArticleReturnsTrueOnSuccess(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("Modify", mock.Anything).Return(CreateSuccessfulModifyResult(), nil)
	result, _ := actions.ArchiveArticle(mockClient, ARTICLE_ID_FIXTURE)
	assert.True(t, result)
}

func TestArchiveArticleReturnsFalseOnClientFailure(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("Modify", mock.Anything).Return(nil, errors.New(MOCK_ERROR_STRING))
	result, _ := actions.ArchiveArticle(mockClient, ARTICLE_ID_FIXTURE)
	assert.False(t, result)
}

func TestArchiveArticleReturnsClientErrorOnClientFailure(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("Modify", mock.Anything).Return(nil, errors.New(MOCK_ERROR_STRING))
	_, err := actions.ArchiveArticle(mockClient, ARTICLE_ID_FIXTURE)
	assert.Equal(t, MOCK_ERROR_STRING, err.Error())
}
