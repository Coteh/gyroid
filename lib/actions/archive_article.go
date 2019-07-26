package actions

import (
	"github.com/Coteh/gyroid/lib/connector"
)

// ArchiveArticle modifies a given Pocket article (by its ID) to archive it from the user's Pocket list
func ArchiveArticle(client connector.PocketConnector, articleID string) (bool, error) {
	return simpleModifyAction(client, articleID, "archive")
}
