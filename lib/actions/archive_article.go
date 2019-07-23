package actions

import (
	pocketConnector "gyroid/lib/connector"
)

// ArchiveArticle modifies a given Pocket article (by its ID) to archive it from the user's Pocket list
func ArchiveArticle(client pocketConnector.PocketConnector, articleID string) (bool, error) {
	return simpleModifyAction(client, articleID, "archive")
}
