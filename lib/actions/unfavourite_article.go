package actions

import (
	pocketConnector "gyroid/lib/connector"
)

// UnfavouriteArticle modifies a given Pocket article (by its ID) to remove it as a favourite
func UnfavouriteArticle(client pocketConnector.PocketConnector, articleID string) (bool, error) {
	return simpleModifyAction(client, articleID, "unfavorite")
}
