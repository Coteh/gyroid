package actions_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/Coteh/gyroid/lib/actions"
	"github.com/Coteh/gyroid/lib/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const ITEM_ID_FIXTURE = "56"
const COUNT_FIXTURE = 200
const START_FIXTURE = 0

type PocketClientUntaggedArticlesMock struct {
	PocketClientMock
}

func (m *PocketClientUntaggedArticlesMock) Retrieve(params models.PocketRetrieve) (*models.PocketRetrieveResult, error) {
	m.Called(params)
	mockArr := make(map[string]models.ArticleResult)
	mockArr[ITEM_ID_FIXTURE] = models.ArticleResult{
		ItemID: ITEM_ID_FIXTURE,
	}

	mockResult := &models.PocketRetrieveResult{
		Status: 0,
		List:   mockArr,
	}

	return mockResult, nil
}

func TestGetUntaggedArticlesGetsUntaggedArticles(t *testing.T) {
	mockClient := &PocketClientUntaggedArticlesMock{}
	articlesList := make([]models.ArticleResult, 0, 20)
	var mut sync.Mutex
	mockClient.On("Retrieve", mock.Anything)
	actions.GetUntaggedArticles(mockClient, 0, 200, &articlesList, &mut)
	assert.Equal(t, ITEM_ID_FIXTURE, articlesList[0].ItemID)
}

func TestGetUnraggedArticlesCallsRetrieveWithCorrectParams(t *testing.T) {
	mockClient := &PocketClientUntaggedArticlesMock{}
	expectedParams := models.PocketRetrieve{
		Tag:    "_untagged_",
		State:  "unread",
		Count:  COUNT_FIXTURE,
		Offset: START_FIXTURE,
		Sort:   "newest",
	}
	articlesList := make([]models.ArticleResult, 0, 20)
	var mut sync.Mutex
	mockClient.On("Retrieve", expectedParams)
	actions.GetUntaggedArticles(mockClient, START_FIXTURE, COUNT_FIXTURE, &articlesList, &mut)
}

func TestGetUntaggedArticlesReturnsClientErrorIfClientError(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("Retrieve", mock.Anything).Return(nil, errors.New(MOCK_ERROR_STRING))
	articlesList := make([]models.ArticleResult, 0, 20)
	var mut sync.Mutex
	err := actions.GetUntaggedArticles(mockClient, 0, 200, &articlesList, &mut)
	assert.Equal(t, MOCK_ERROR_STRING, err.Error())
}

func TestGetUntaggedArticlesShouldNotAddAnyArticlesIfClientError(t *testing.T) {
	mockClient := &PocketClientMock{}
	mockClient.On("Retrieve", mock.Anything).Return(nil, errors.New(MOCK_ERROR_STRING))
	articlesList := make([]models.ArticleResult, 0, 20)
	var mut sync.Mutex
	actions.GetUntaggedArticles(mockClient, 0, 200, &articlesList, &mut)
	assert.Equal(t, len(articlesList), 0)
}

func TestGetUntaggedArticlesShouldNotAddAnyArticlesIfStartIsInvalid(t *testing.T) {
	mockClient := &PocketClientUntaggedArticlesMock{}
	articlesList := make([]models.ArticleResult, 0, 20)
	var mut sync.Mutex
	mockClient.On("Retrieve", mock.Anything)
	actions.GetUntaggedArticles(mockClient, -1, 200, &articlesList, &mut)
	assert.Zero(t, len(articlesList))
}

func TestGetUntaggedArticlesShouldReturnErrorIfStartIsInvalid(t *testing.T) {
	mockClient := &PocketClientUntaggedArticlesMock{}
	articlesList := make([]models.ArticleResult, 0, 20)
	var mut sync.Mutex
	mockClient.On("Retrieve", mock.Anything)
	err := actions.GetUntaggedArticles(mockClient, -1, 200, &articlesList, &mut)
	assert.NotNil(t, err)
}

func TestGetUntaggedArticlesShouldNotMakeAnyRequestIfStartIsInvalid(t *testing.T) {
	mockClient := &PocketClientUntaggedArticlesMock{}
	articlesList := make([]models.ArticleResult, 0, 20)
	var mut sync.Mutex
	mockClient.On("Retrieve", mock.Anything)
	actions.GetUntaggedArticles(mockClient, -1, 200, &articlesList, &mut)
	mockClient.AssertNumberOfCalls(t, "Retrieve", 0)
}

func TestGetUntaggedArticlesShouldNotAddAnyArticlesIfCountIsZero(t *testing.T) {
	mockClient := &PocketClientUntaggedArticlesMock{}
	articlesList := make([]models.ArticleResult, 0, 20)
	var mut sync.Mutex
	mockClient.On("Retrieve", mock.Anything)
	actions.GetUntaggedArticles(mockClient, 0, 0, &articlesList, &mut)
	assert.Zero(t, len(articlesList))
}

func TestGetUntaggedArticlesShouldReturnErrorIfCountIsZero(t *testing.T) {
	mockClient := &PocketClientUntaggedArticlesMock{}
	articlesList := make([]models.ArticleResult, 0, 20)
	var mut sync.Mutex
	mockClient.On("Retrieve", mock.Anything)
	err := actions.GetUntaggedArticles(mockClient, 0, 0, &articlesList, &mut)
	assert.NotNil(t, err)
}

func TestGetUntaggedArticlesShouldNotMakeAnyRequestIfCountIsZero(t *testing.T) {
	mockClient := &PocketClientUntaggedArticlesMock{}
	articlesList := make([]models.ArticleResult, 0, 20)
	var mut sync.Mutex
	mockClient.On("Retrieve", mock.Anything)
	actions.GetUntaggedArticles(mockClient, 0, 0, &articlesList, &mut)
	mockClient.AssertNumberOfCalls(t, "Retrieve", 0)
}
