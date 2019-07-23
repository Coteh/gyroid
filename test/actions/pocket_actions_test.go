package actions_test

import (
	"errors"
	"gyroid/lib/models"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const ARTICLE_ID_FIXTURE = "100"
const MOCK_ERROR_STRING = "Test"

// PocketClientMock implements PocketConnector interface but is not used directly
type PocketClientMock struct {
	mock.Mock
}

func (m *PocketClientMock) SetAccessToken(accessToken string) {

}

func (m *PocketClientMock) Retrieve(params models.PocketRetrieve) (*models.PocketRetrieveResult, error) {
	return nil, nil
}

func (m *PocketClientMock) Add(params models.PocketAdd) (*models.PocketAddResult, error) {
	return nil, nil
}

func (m *PocketClientMock) Modify(params models.PocketModify) (*models.PocketModifyResult, error) {
	return nil, nil
}

func (m *PocketClientMock) RequestOAuthCode(redirectUri string) (string, error) {
	return "", nil
}

func (m *PocketClientMock) Authorize(code string) (string, error) {
	return "", nil
}

type FailingPocketClientMock struct {
}

func (m *FailingPocketClientMock) SetAccessToken(accessToken string) {

}

func (m *FailingPocketClientMock) Retrieve(params models.PocketRetrieve) (*models.PocketRetrieveResult, error) {
	return nil, errors.New(MOCK_ERROR_STRING)
}

func (m *FailingPocketClientMock) Add(params models.PocketAdd) (*models.PocketAddResult, error) {
	return nil, errors.New(MOCK_ERROR_STRING)
}

func (m *FailingPocketClientMock) Modify(params models.PocketModify) (*models.PocketModifyResult, error) {
	return nil, errors.New(MOCK_ERROR_STRING)
}

func (m *FailingPocketClientMock) RequestOAuthCode(redirectUri string) (string, error) {
	return "", errors.New(MOCK_ERROR_STRING)
}

func (m *FailingPocketClientMock) Authorize(code string) (string, error) {
	return "", errors.New(MOCK_ERROR_STRING)
}
