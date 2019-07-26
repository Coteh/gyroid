package actions

import (
	"errors"
	pocketConnector "github.com/Coteh/gyroid/lib/connector"
	models "github.com/Coteh/gyroid/lib/models"

	"sync"
)

// GetUntaggedArticles retrieves all untagged articles for a given user (by their access token) from Pocket
func GetUntaggedArticles(client pocketConnector.PocketConnector, start int, count int, articleList *[]models.ArticleResult, mut *sync.Mutex) error {
	if start < 0 {
		return errors.New("Start is invalid")
	} else if count == 0 {
		return errors.New("Count should not be 0")
	}

	params := models.PocketRetrieve{
		Tag:    "_untagged_",
		State:  "unread",
		Count:  count,
		Offset: start,
		Sort:   "newest",
	}

	resp, err := client.Retrieve(params)
	if err != nil {
		return err
	}

	articlesMap := resp.List

	mut.Lock()
	defer mut.Unlock()

	for _, article := range articlesMap {
		*articleList = append(*articleList, article)
	}

	return nil
}
