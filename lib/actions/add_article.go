package actions

import (
	pocketConnector "gyroid/lib/connector"
	models "gyroid/lib/models"
)

// AddArticle adds a given article (by its URL) to the user's Pocket list
func AddArticle(client pocketConnector.PocketConnector, url string) (map[string]interface{}, error) {
	addParams := models.PocketAdd{
		Url: url,
	}

	result, err := client.Add(addParams)
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}
