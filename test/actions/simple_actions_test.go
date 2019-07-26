package actions_test

import (
	"github.com/Coteh/gyroid/lib/models"
)

type PocketClientSimpleActionMock struct {
	PocketClientMock
}

func (m *PocketClientSimpleActionMock) Modify(params models.PocketModify) (*models.PocketModifyResult, error) {
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
