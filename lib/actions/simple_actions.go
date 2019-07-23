package actions

import (
	pocketConnector "gyroid/lib/connector"
	models "gyroid/lib/models"
)

func simpleModifyAction(client pocketConnector.PocketConnector, articleID string, action string) (bool, error) {
	resultArr, err := sendModifyRequest(client, &models.PocketAction{
		Action: action,
		ItemID: articleID,
	})

	if err != nil {
		return false, err
	}

	return resultArr[0].(bool), nil
}
