package models

// PocketParam is an interface representing a request parameter to Pocket API
type PocketParam interface{}

// PocketResponse is an interface representing a response from Pocket API
type PocketResponse interface{}

// PocketResult is an interface representing a result item from Pocket API
type PocketResult interface{}

// PocketAdd represents a request to Pocket API Add endpoint
type PocketAdd struct {
	Url     string `json:"url"`
	Title   string `json:"title,omitempty"`
	Tags    string `json:"tags,omitempty"`
	TweetID string `json:"tweet_id,omitempty"`
}

// PocketAction represents a single action item for the Pocket API Modify endpoint
type PocketAction struct {
	Action string `json:"action"`
	ItemID string `json:"item_id"`
	Time   string `json:"time,omitempty"`
	Tags   string `json:"tags,omitempty"`
}

// PocketModify represents a request to Pocket API Modify endpoint
type PocketModify struct {
	Actions []PocketAction `json:"actions"`
}

// PocketRetrieve represents a request to Pocket API Retrieve endpoint
type PocketRetrieve struct {
	Tag    string `json:"tag"`
	Sort   string `json:"sort"`
	State  string `json:"state"`
	Count  int    `json:"count"`
	Offset int    `json:"offset"`
}

// PocketAddResult represents results from Pocket API's Add endpoint
// TODO create a struct type for Item
type PocketAddResult struct {
	Status int                    `json:"status"`
	Item   map[string]interface{} `json:"item"`
}

// PocketModifyResult represents results from Pocket API's Modify endpoint
// TODO create two struct types for action results and errors
type PocketModifyResult struct {
	Status        int           `json:"status"`
	ActionResults []interface{} `json:"action_results"`
	ActionErrors  []interface{} `json:"action_errors"`
}

// PocketRetrieveResult represents results from Pocket API's Retrieve endpoint
type PocketRetrieveResult struct {
	Status int                      `json:"status"`
	List   map[string]ArticleResult `json:"list"`
}
