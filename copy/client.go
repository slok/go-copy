// Here starts all the mess :)
//
// You need a client to work with all the time, the client has all the neccessary
// things to work with: URL, Session, result decoding logic...
//
// First create a client and the create services with the created client, these
// services will be the ones that retrieve data from the Copy servers
//
//
// The program has some global package variables
//
// defaultHttpClient: The default http client
// appTokenEnv: The copy app oauth token
// appSecretEnv: The copy app oauth secret
// accessTokenEnv : The user authorized oauth token for the app
// accessSecretEnv: The user authorized oauth secret for the app
// session: The session for the oauth hand shaking
// mux: the mux for the server mocking in the tests
// client: The mighty client for the job ;)
// server: The mock server for the tests

package copy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
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

// Oauth handshake neccesary data
const (
	appTokenEnv     = "APP_TOKEN"
	appSecretEnv    = "APP_SECRET"
	accessTokenEnv  = "ACCESS_TOKEN"
	accessSecretEnv = "ACCESS_SECRET"
)

// Global vars for the tokens
var (
	appToken     string
	appSecret    string
	accessToken  string
	accessSecret string
)

// Creates a new client. If no http client and URL the client will use the
// default ones
func NewClient(httpClient *http.Client, resourcesUrl string,
	appToken string, appSecret string,
	accessToken string, accessSecret string) (*Client, error) {

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
	}

	c.session = session
	return c, nil

}

// Returns a default client, normally we will use this
func NewDefaultClient(appToken string, appSecret string,
	accessToken string, accessSecret string) (*Client, error) {
	return NewClient(nil, "", appToken, appSecret, accessToken, accessSecret)
}

// Makes the client request based on the url, method, values and returns
// the response is the response of the call
// the value is inside the v param (you should pass a pointer because will
// mutate inside the method)
func (c *Client) DoRequestDecoding(method string, urlStr string, form url.Values, v interface{}) (*http.Response, error) {
	var resp *http.Response
	var err error

	endpoint := strings.Join([]string{c.resourcesUrl, urlStr}, "/")

	switch method {
	case "GET":
		resp, err = c.session.Get(endpoint, form, c.httpClient)

	case "POST":
		//resp, err = c.session.Post(endpoint, form, c.httpClient)

	case "PUT":
		resp, err = c.session.Put(endpoint, form, c.httpClient)

	case "DELETE":
		resp, err = c.session.Delete(endpoint, form, c.httpClient)

	}

	if err != nil || resp == nil {
		return nil, errors.New("Error making the request")
	}

	defer resp.Body.Close()

	// If v is nil that means that the caller doesn't need the response
	if v != nil {
		// response body to string
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)

		if err != nil {
			return nil, err
		}
		respBody := buf.String()

		// Decode to our structure
		json.NewDecoder(strings.NewReader(respBody)).Decode(v)
	}

	return resp, nil
}

// Makes the client request based on the url.
//
// This will be binary data body so we don't process the request
func (c *Client) DoRequestContent(urlStr string) (*http.Response, error) {
	var resp *http.Response
	var err error

	endpoint := strings.Join([]string{c.resourcesUrl, urlStr}, "/")

	resp, err = c.session.Get(endpoint, nil, c.httpClient)

	if err != nil || resp == nil {
		return nil, errors.New("Error making the request")
	}

	// Don't close the body is a chunked HTTP response
	//defer resp.Body.Close()

	return resp, nil
}

// Makes the client request for uploading multipart request
//
func (c *Client) DoRequestMultipart(filePath, uploadPath, filename string) (*http.Response, error) {

	endpoint := strings.Join([]string{c.resourcesUrl, uploadPath}, "/")

	// Do sequential upload (not all in memory)

	// Get our file reader
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	// The key is to use a pipe
	// (reading the file will feed the writer and then the reader will be in the multipart body
	// so when the multipart request a chunk the reader will "call" the writer and this will
	// read from the file)
	//
	// File -> FileReader -> Writer -> (Multipart wrapp magic) -> Reader -> Multipart body
	/*reader, writer := io.Pipe()
	multiWriter := multipart.NewWriter(writer)

	// The sequential write (read from file) will be in a goroutine
	go func() {
		defer writer.Close()
		defer file.Close()

		part, _ := multiWriter.CreateFormFile("file", filename)

		// Copy on demand
		io.Copy(part, file)
		multiWriter.Close()
	}()

	// This will be custom because the multipart is trickier thatn a normal request
	req, err := http.NewRequest("POST", endpoint, reader)
	if err != nil {
		return nil, err
	}*/

	//-----------------------------------------------------------
	// FIXME: See above, use pipes to read/write and don't load in memory
	// From the docs: The maximum filesize of an upload is 1GB. An API endpoint
	// supporting chunked file uploading is planned for circumventing this limitation.
	//defer file.Close()

	body := &bytes.Buffer{}
	multiWriter := multipart.NewWriter(body)
	part, err := multiWriter.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	multiWriter.Close()
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, err
	}
	//-----------------------------------------------------------

	fileInfo, _ := file.Stat()
	req.Header.Set("Authorization", c.session.OauthClient.AuthorizationHeader(&c.session.TokenCreds, "POST", req.URL, nil))
	req.Header.Set("Content-Type", multiWriter.FormDataContentType())
	req.Header.Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	return c.session.Do(req, c.httpClient)
}
