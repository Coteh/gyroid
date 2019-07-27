package actions

import (
	"github.com/Coteh/gyroid/lib/connector"
	"github.com/Coteh/gyroid/lib/models"
)

func simpleModifyAction(client connector.PocketConnector, articleID string, action string) (bool, error) {
	resultArr, err := sendModifyRequest(client, &models.PocketAction{
		Action: action,
		ItemID: articleID,
	})

	if err != nil {
		return false, err
	}

	return resultArr[0].(bool), nil
}
