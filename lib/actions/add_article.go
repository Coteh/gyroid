package actions

import (
	"github.com/Coteh/gyroid/lib/connector"
	"github.com/Coteh/gyroid/lib/models"
)

// AddArticle adds a given article (by its URL) to the user's Pocket list
func AddArticle(client connector.PocketConnector, url string) (*models.AddedArticleResult, error) {
	addParams := models.PocketAdd{
		Url: url,
	}

	result, err := client.Add(addParams)
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}
