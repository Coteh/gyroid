package actions

import (
	"github.com/Coteh/gyroid/lib/connector"
)

// DeleteArticle modifies a given Pocket article (by its ID) to delete it from the user's Pocket list
func DeleteArticle(client connector.PocketConnector, articleID string) (bool, error) {
	return simpleModifyAction(client, articleID, "delete")
}
