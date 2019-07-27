package actions

import (
	"github.com/Coteh/gyroid/lib/connector"
)

// UnfavouriteArticle modifies a given Pocket article (by its ID) to remove it as a favourite
func UnfavouriteArticle(client connector.PocketConnector, articleID string) (bool, error) {
	return simpleModifyAction(client, articleID, "unfavorite")
}
