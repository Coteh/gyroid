package actions_test

import (
	"github.com/Coteh/gyroid/lib/models"

	"errors"

	"github.com/stretchr/testify/mock"
)

const ARTICLE_ID_FIXTURE = "100"
const MOCK_ERROR_STRING = "Test"

// PocketClientMock implements PocketConnector interface and used for tests
type PocketClientMock struct {
	mock.Mock
}

func (m *PocketClientMock) SetAccessToken(accessToken string) {

}

func (m *PocketClientMock) Retrieve(params models.PocketRetrieve) (*models.PocketRetrieveResult, error) {
	args := m.Called(params)
	var result *models.PocketRetrieveResult
	if args.Get(0) != nil {
		result = args.Get(0).(*models.PocketRetrieveResult)
	}
	return result, args.Error(1)
}

func (m *PocketClientMock) Add(params models.PocketAdd) (*models.PocketAddResult, error) {
	args := m.Called(params)
	var result *models.PocketAddResult
	if args.Get(0) != nil {
		result = args.Get(0).(*models.PocketAddResult)
	}
	return result, args.Error(1)
}

func (m *PocketClientMock) Modify(params models.PocketModify) (*models.PocketModifyResult, error) {
	args := m.Called(params)
	var result *models.PocketModifyResult
	if args.Get(0) != nil {
		result = args.Get(0).(*models.PocketModifyResult)
	}
	return result, args.Error(1)
}

func (m *PocketClientMock) RequestOAuthCode(redirectUri string) (string, error) {
	args := m.Called(redirectUri)
	return args.String(0), args.Error(1)
}

func (m *PocketClientMock) Authorize(code string) (string, error) {
	args := m.Called(code)
	return args.String(0), args.Error(1)
}

func CreateSuccessfulRetrieveResult(itemID string) *models.PocketRetrieveResult {
	mockArr := make(map[string]models.ArticleResult)
	mockArr[itemID] = models.ArticleResult{
		ItemID: itemID,
	}
	return &models.PocketRetrieveResult{
		Status: 0,
		List:   mockArr,
	}
}

func CreateSuccessfulAddResult(url string) *models.PocketAddResult {
	mockItem := make(map[string]interface{})
	mockItem["normal_url"] = url
	return &models.PocketAddResult{
		Status: 0,
		Item:   mockItem,
	}
}

func CreateSuccessfulModifyResult() *models.PocketModifyResult {
	mockArr := make([]interface{}, 1)
	mockArr[0] = true
	return &models.PocketModifyResult{
		Status:        0,
		ActionResults: mockArr,
		ActionErrors:  make([]interface{}, 0),
	}
}

func CreatePocketClientMock() *PocketClientMock {
	return &PocketClientMock{}
}

func CreatePocketClientMockWithExpectation(methodName string, arguments ...interface{}) (*PocketClientMock, *mock.Call) {
	mockClient := &PocketClientMock{}
	return mockClient, mockClient.On(methodName, arguments...)
}

func CreateFailingPocketClientMock(methodName string, arguments ...interface{}) *PocketClientMock {
	mockClient, call := CreatePocketClientMockWithExpectation(methodName, arguments...)
	call.Return(nil, errors.New(MOCK_ERROR_STRING))
	return mockClient
}

func CreateSuccessfulPocketClientMock(methodName string, result interface{}, arguments ...interface{}) *PocketClientMock {
	mockClient, call := CreatePocketClientMockWithExpectation(methodName, arguments...)
	call.Return(result, nil)
	return mockClient
}
