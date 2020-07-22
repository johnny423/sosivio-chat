package auth

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//NewRESTClient creates a new instance of RESTClient from using the host's
//address and a user-agent for the client
func NewRESTClient(host string, userAgent string, logger log.FieldLogger) *RESTClient {
	clientUrl, _ := url.Parse(host)
	return &RESTClient{
		BaseURL:    clientUrl,
		UserAgent:  userAgent,
		httpClient: &http.Client{Timeout: time.Minute},
		logger:     logger,
	}
}

//RESTClient is a REST API client for auth server
type RESTClient struct {
	BaseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
	logger     log.FieldLogger
}

//Validate checks with the remote server if the token is validate and if so, it returns the token's username
func (client *RESTClient) Validate(token string) (string, error) {
	resp, err := client.sendRequest("POST", &TokenMsg{Token: token}, "/validate")
	if err != nil {
		client.logger.Error("error raised while sending a request ", err)
		return "", err
	}
	defer resp.Body.Close()

	var result ValidateResult

	body, _ := ioutil.ReadAll(resp.Body)
	client.logger.Debug("the respond body is:", string(body))

	err = json.Unmarshal(body, &result)
	if err != nil {
		client.logger.Error("unmarshal ", err)
		return "", err
	}
	client.logger.Debug("result:", result)
	return result.UserName, nil
}

//sendRequest sends a new request from the given method, request body and resource and returns the response
func (client *RESTClient) sendRequest(method string, requestValue interface{}, resource string) (*http.Response, error) {
	u := client.resolveUrl(resource)
	v, err := json.Marshal(requestValue)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(v))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", client.UserAgent)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//resolveUrl creates a url from the client base url and a given resource
func (client *RESTClient) resolveUrl(resource string) *url.URL {
	rel := &url.URL{Path: resource}
	u := client.BaseURL.ResolveReference(rel)
	return u
}
