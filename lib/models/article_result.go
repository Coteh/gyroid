package models

// ArticleResult represents all common results
type ArticleResult struct {
	ItemID        string `json:"item_id"`
	ResolvedID    string `json:"resolved_id"`
	ResolvedURL   string `json:"resolved_url"`
	ResolvedTitle string `json:"resolved_title"`

	GivenURL   string `json:"given_url"`
	GivenTitle string `json:"given_title"`
	Favorite   int    `json:"favorite,string"`
	Status     int    `json:"status,string"`

	Excerpt   string `json:"excerpt"`
	IsArticle string `json:"is_article"`
	HasVideo  int    `json:"has_video,string"`
	HasImage  int    `json:"has_image,string"`
	WordCount int    `json:"word_count,string"`
	// Authors   []string `json:"authors,string"`
	// Images    []string `json:"images,string"`
	// Videos    []string `json:"videos,string"`
}

// AddedArticleResult represents all results for add article
type AddedArticleResult struct {
	*ArticleResult
	NormalURL      string `json:"normal_url"`
	DomainID       string `json:"domain_id"`
	OriginDomainID string `json:"origin_domain_id"`
	ResponseCode   string `json:"response_code"`
	MimeType       string `json:"mime_type"`
	ContentLength  string `json:"content_length"`
	Encoding       string `json:"encoding"`
	DateResolved   string `json:"date_resolved"`
	DatePublished  string `json:"date_published"`
	Title          string `json:"title"`
	IsIndex        string `json:"isindex"`
}

// RetrievedArticleResult represents all results for retrieve article
type RetrievedArticleResult struct {
	*ArticleResult
	Tags []string `json:"tags,string"`
}
