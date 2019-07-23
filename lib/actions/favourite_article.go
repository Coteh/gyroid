package actions

import (
	pocketConnector "gyroid/lib/connector"
)

// FavouriteArticle modifies a given Pocket article (by its ID) to mark it as a favourite
func FavouriteArticle(client pocketConnector.PocketConnector, articleID string) (bool, error) {
	return simpleModifyAction(client, articleID, "favorite")
}
