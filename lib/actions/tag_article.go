package actions

import (
	"github.com/Coteh/gyroid/lib/connector"
	"github.com/Coteh/gyroid/lib/models"

	"strings"
)

// MarkArticleWithTag modifies a given Pocket article (by its ID) to give it a specified set of tags
func MarkArticleWithTag(client connector.PocketConnector, articleID string, tags []string) (bool, error) {
	resultArr, err := sendModifyRequest(client, &models.PocketAction{
		Action: "tags_add",
		ItemID: articleID,
		Tags:   ConcatTags(tags),
	})

	if err != nil {
		return false, err
	}

	return resultArr[0].(bool), nil
}

// ConcatTags concatenates tags into a format suitable for Pocket requests involving tags
func ConcatTags(tags []string) string {
	return strings.Join(tags, ",")
}
