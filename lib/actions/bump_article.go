package actions

import (
	pocketConnector "gyroid/lib/connector"
	models "gyroid/lib/models"
)

// BumpArticleToTop modifies a given Pocket article (by its ID) to readd (unarchive) it from the user's list
// this also acts as a way of "bumping" the article back up to the top of the user's Pocket list
func BumpArticleToTop(client pocketConnector.PocketConnector, articleID string) (bool, error) {
	resultArr, err := sendModifyRequest(client, &models.PocketAction{
		Action: "readd",
		ItemID: articleID,
	})

	if err != nil {
		return false, err
	}

	return resultArr[0] != nil, nil
}
