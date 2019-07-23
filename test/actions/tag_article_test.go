package actions_test

import (
	"gyroid/lib/actions"
	"gyroid/lib/models"
	"testing"
	// "github.com/stretchr/testify/assert"
)

type PocketClientMarkArticleWithTagMock struct {
	PocketClientMock
}

func (m *PocketClientMarkArticleWithTagMock) Modify(params models.PocketModify) (*models.PocketModifyResult, error) {
	m.Called()
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
	t.Fatal("Not implemented - waiting on param checking tests till all other tests are done, then will install testify and mockery")
}

func TestMarkArticleWithTagReturnsTrueOnSuccess(t *testing.T) {
	mockClient := &PocketClientMarkArticleWithTagMock{}
	tags := make([]string, 1)
	tags[0] = "test"
	result, _ := actions.MarkArticleWithTag(mockClient, "100", tags)
	if !result {
		t.Logf("Expected true from MarkArticleWithTag - got false")
		t.FailNow()
	}
}

func TestMarkArticleWithTagReturnsFalseOnClientFailure(t *testing.T) {
	mockClient := &FailingPocketClientMock{}
	tags := make([]string, 1)
	tags[0] = "test"
	result, _ := actions.MarkArticleWithTag(mockClient, "100", tags)
	if result {
		t.Logf("Expected false from MarkArticleWithTag - got true")
		t.FailNow()
	}
}

func TestMarkArticleWithTagReturnsClientErrorOnClientFailure(t *testing.T) {
	mockClient := &FailingPocketClientMock{}
	tags := make([]string, 1)
	tags[0] = "test"
	_, err := actions.MarkArticleWithTag(mockClient, "100", tags)
	if err == nil || err.Error() != MOCK_ERROR_STRING {
		t.Logf("Did not receive client error")
		t.FailNow()
	}
}
