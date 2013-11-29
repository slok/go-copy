package copy

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Client has the session (Session) for calling the REST API with Oauth
// the Http client and the URL to call
type Client struct {
	session      *Session
	resourcesUrl string
	httpClient   *http.Client
}

const (
	defaultResourcesUrl = "https://api.copy.com/rest"
)

var (
	defaultHttpClient = http.DefaultClient
)

// Creates a new client. If no http client and URL the client will use the
// default ones
func NewClient(httpClient *http.Client, resourcesUrl string) (*Client, error) {

	c := new(Client)

	if httpClient == nil {
		c.httpClient = defaultHttpClient
	} else {
		c.httpClient = httpClient
	}

	if resourcesUrl == "" {
		c.resourcesUrl = defaultResourcesUrl
	} else {
		c.resourcesUrl = resourcesUrl
	}

	// Create session
	//FIX: For now use Env vars, next params and/or conf file
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = os.Getenv(accessSecretEnv)

	session, err := NewSession(
		AppToken{
			Token: appToken,
			Key:   appSecret,
		},
		AccessToken{
			Token: accessToken,
			Key:   accessSecret,
		},
	)

	if err != nil || appToken == "" || appSecret == "" || accessToken == "" ||
		accessSecret == "" {
		return nil, errors.New("Could not create the client, Check access settings")
	} else {
		c.session = session
		return c, nil
	}
}

// Makes the client request based on the url, method, values and returns
// the response is the response of the call
// the value is inside the t param (you should pass a pointer because will
// mutate inside the method)
func (c *Client) Do(method string, urlStr string, form url.Values, v interface{}) (*http.Response, error) {
	var resp *http.Response
	var err error

	endpoint := strings.Join([]string{c.resourcesUrl, urlStr}, "/")

	switch method {
	case "GET":
		resp, err = c.session.Get(endpoint, form, c.httpClient)

	case "POST":
		resp, err = c.session.Post(endpoint, form, c.httpClient)

	case "PUT":
		resp, err = c.session.Put(endpoint, form, c.httpClient)

	case "DELETE":
		resp, err = c.session.Delete(endpoint, form, c.httpClient)

	}

	if err != nil || resp == nil {
		return nil, err
	}

	defer resp.Body.Close()

	// response body to string
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)

	if err != nil {
		return nil, err
	}
	respBody := buf.String()

	// Decode to our structure
	json.NewDecoder(strings.NewReader(respBody)).Decode(v)

	return resp, nil
}