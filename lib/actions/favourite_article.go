package actions

import (
	"github.com/Coteh/gyroid/lib/connector"
)

// FavouriteArticle modifies a given Pocket article (by its ID) to mark it as a favourite
func FavouriteArticle(client connector.PocketConnector, articleID string) (bool, error) {
	return simpleModifyAction(client, articleID, "favorite")
}
