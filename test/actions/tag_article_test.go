package actions_test

import (
	"github.com/Coteh/gyroid/lib/actions"
	"github.com/Coteh/gyroid/lib/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type PocketClientMarkArticleWithTagMock struct {
	PocketClientMock
}

func (m *PocketClientMarkArticleWithTagMock) Modify(params models.PocketModify) (*models.PocketModifyResult, error) {
	m.Called(params)
	mockArr := make([]interface{}, 1)
	mockArr[0] = true
	mockResult := &models.PocketModifyResult{
		Status:        0,
		ActionResults: mockArr,
		ActionErrors:  make([]interface{}, 0),
	}

	return mockResult, nil
}

func TestMarkArticleWithTagCallsModifyWithCorrectParams(t *testing.T) {
	mockClient := &PocketClientMarkArticleWithTagMock{}
	tags := make([]string, 1)
	tags[0] = "test"
	expectedTagReqStr := tags[0]
	expectedAction := models.PocketAction{
		Action: "tags_add",
		ItemID: ARTICLE_ID_FIXTURE,
		Tags:   expectedTagReqStr,
	}
	expectedActionArr := make([]models.PocketAction, 1)
	expectedActionArr[0] = expectedAction
	expectedParams := models.PocketModify{
		Actions: expectedActionArr,
	}
	mockClient.On("Modify", expectedParams)
	actions.MarkArticleWithTag(mockClient, ARTICLE_ID_FIXTURE, tags)
}

func TestMarkArticleWithTagReturnsTrueOnSuccess(t *testing.T) {
	mockClient := &PocketClientMarkArticleWithTagMock{}
	tags := make([]string, 1)
	tags[0] = "test"
	mockClient.On("Modify", mock.Anything)
	result, _ := actions.MarkArticleWithTag(mockClient, "100", tags)
	assert.True(t, result)
}

func TestMarkArticleWithTagReturnsFalseOnClientFailure(t *testing.T) {
	mockClient := &FailingPocketClientMock{}
	tags := make([]string, 1)
	tags[0] = "test"
	result, _ := actions.MarkArticleWithTag(mockClient, "100", tags)
	assert.False(t, result)
}

func TestMarkArticleWithTagReturnsClientErrorOnClientFailure(t *testing.T) {
	mockClient := &FailingPocketClientMock{}
	tags := make([]string, 1)
	tags[0] = "test"
	_, err := actions.MarkArticleWithTag(mockClient, "100", tags)
	assert.Equal(t, MOCK_ERROR_STRING, err.Error())
}

func TestConcatTagsConcatenatesTags(t *testing.T) {
	tags := make([]string, 2)
	tags[0] = "cat"
	tags[1] = "dog"
	assert.Equal(t, "cat,dog", actions.ConcatTags(tags))
}

func TestConcatTagsDoesNotAlterSingleTag(t *testing.T) {
	tags := make([]string, 1)
	tags[0] = "sameness"
	assert.Equal(t, "sameness", actions.ConcatTags(tags))
}
