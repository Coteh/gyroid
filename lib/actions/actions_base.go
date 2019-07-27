package actions

import (
	"github.com/Coteh/gyroid/lib/connector"
	"github.com/Coteh/gyroid/lib/models"
)

func sendModifyRequest(client connector.PocketConnector, payload *models.PocketAction) ([]interface{}, error) {
	actions := make([]models.PocketAction, 1)
	actions[0] = *payload

	params := models.PocketModify{
		Actions: actions,
	}

	result, err := client.Modify(params)
	if err != nil {
		return nil, err
	}

	return result.ActionResults, nil
}
