// Package connector contains functions that interface with Pocket's API
package connector

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Coteh/gyroid/lib/models"
	"io/ioutil"
	"net/http"
)

// PocketConnector is an interface for methods that call Pocket API
type PocketConnector interface {
	SetAccessToken(accessToken string)
	Retrieve(params models.PocketRetrieve) (*models.PocketRetrieveResult, error)
	Add(params models.PocketAdd) (*models.PocketAddResult, error)
	Modify(params models.PocketModify) (*models.PocketModifyResult, error)
	RequestOAuthCode(redirectURI string) (string, error)
	Authorize(code string) (string, error)
}

// PocketClient represents a client connection to Pocket
type PocketClient struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
	pocketURL   string
}

// GetJSON combines given request params and Pocket client information into a single JSON byte array
func GetJSON(params models.PocketParam, pocketClient *PocketClient) ([]byte, error) {
	mappedParams, err := createMapFromParams(params)
	if err != nil {
		return nil, err
	}

	clientParams, err := createMapFromParams(pocketClient)
	if err != nil {
		return nil, err
	}

	for key, val := range mappedParams {
		clientParams[key] = val
	}

	result, err := json.Marshal(clientParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreatePocketClient creates a PocketClient instance that can be used to
// interface with Pocket API functions
func CreatePocketClient(consumerKey string, accessToken string) *PocketClient {
	return &PocketClient{
		ConsumerKey: consumerKey,
		AccessToken: accessToken,
		pocketURL:   "https://getpocket.com/v3",
	}
}

func (client *PocketClient) getPocketURL(subPath string) string {
	return client.pocketURL + subPath
}

// SetAccessToken sets the access token on a PocketClient instance
func (client *PocketClient) SetAccessToken(accessToken string) {
	client.AccessToken = accessToken
}

// Retrieve calls Pocket API's Retrieve method given arguments in params
func (client *PocketClient) Retrieve(params models.PocketRetrieve) (*models.PocketRetrieveResult, error) {
	jsonBytes, err := GetJSON(params, client)
	if err != nil {
		return nil, err
	}

	r := bytes.NewBuffer(jsonBytes)

	httpClient := http.Client{}
	var req *http.Request
	var resp *http.Response
	var body []byte

	req, err = http.NewRequest("POST", client.getPocketURL("/get"), r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Accept", "application/json")

	resp, err = httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := new(models.PocketRetrieveResult)

	err = json.Unmarshal(body, &result)
	if err != nil {
		// Status of 2 means no new articles, just return result with empty list
		if result.Status == 2 {
			return result, nil
		}
		return nil, err
	}

	return result, nil
}

// Add calls Pocket API's Add method to add the URL in parameter to the
// Pocket list of the user associated with access token
func (client *PocketClient) Add(params models.PocketAdd) (*models.PocketAddResult, error) {
	jsonBytes, err := GetJSON(params, client)
	if err != nil {
		return nil, err
	}

	r := bytes.NewBuffer(jsonBytes)

	httpClient := http.Client{}
	var req *http.Request
	var resp *http.Response
	var body []byte

	req, err = http.NewRequest("POST", client.getPocketURL("/add"), r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Accept", "application/json")

	resp, err = httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := new(models.PocketAddResult)

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Modify calls Pocket API's Modify method with provided arguments
// from actions parameter
func (client *PocketClient) Modify(params models.PocketModify) (*models.PocketModifyResult, error) {
	jsonBytes, err := GetJSON(params, client)
	if err != nil {
		return nil, err
	}

	r := bytes.NewBuffer(jsonBytes)

	httpClient := http.Client{}
	var req *http.Request
	var resp *http.Response
	var body []byte

	req, err = http.NewRequest("POST", client.getPocketURL("/send"), r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Accept", "application/json")

	resp, err = httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pocketResp := new(models.PocketModifyResult)

	err = json.Unmarshal(body, &pocketResp)
	if err != nil {
		return nil, err
	}

	if pocketResp.ActionResults == nil {
		return nil, errors.New("Action results for Pocket modify action not provided")
	}

	return pocketResp, nil
}

// RequestOAuthCode requests an OAuth request token from Pocket's API to provide to a user
func (client *PocketClient) RequestOAuthCode(redirectURI string) (string, error) {
	reqParams := make(map[string]string)
	httpClient := http.Client{}
	var r *bytes.Buffer

	reqParams["consumer_key"] = client.ConsumerKey
	reqParams["redirect_uri"] = redirectURI

	b, err := json.Marshal(reqParams)
	if err != nil {
		return "", err
	}

	r = bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", client.getPocketURL("/oauth/request"), r)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	oauthReqResp := make(map[string]string)

	err = json.Unmarshal(body, &oauthReqResp)
	if err != nil {
		return "", err
	}

	code := oauthReqResp["code"]

	return code, nil
}

// Authorize authorizes a user associated with an OAuth request token
// with the Pocket API and returns an access token that can then be used to make calls
// to the Pocket API for that user
func (client *PocketClient) Authorize(code string) (string, error) {
	reqParams := make(map[string]string)
	httpClient := http.Client{}
	reqParams["code"] = code
	reqParams["consumer_key"] = client.ConsumerKey

	b, err := json.Marshal(reqParams)
	if err != nil {
		return "", err
	}

	r := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", client.getPocketURL("/oauth/authorize"), r)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Pocket Error occurred - %s", resp.Header.Get("X-Error-Code"))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	accessResult := make(map[string]string)

	err = json.Unmarshal(body, &accessResult)
	if err != nil {
		return "", err
	}

	accessToken := accessResult["access_token"]

	return accessToken, nil
}

func createMapFromParams(paramsObj interface{}) (map[string]interface{}, error) {
	var mappedParams map[string]interface{}

	paramsStr, err := json.Marshal(paramsObj)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(paramsStr, &mappedParams)

	return mappedParams, err
}
