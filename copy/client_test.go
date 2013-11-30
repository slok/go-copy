package copy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

// Global testing vars
var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

// Setups the mock server for all the app tests
func setup(t *testing.T) {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	var err error
	client, err = NewTestClient()
	if err != nil {
		t.Error(err.Error())
	}
}

// ---------Global testing utils-----------------------------------------------

// Cleans all the setup of the test
func tearDown() {
	server.Close()
}

// From go-github (https://github.com/google/go-github)
func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

// Creates a new Client for testing
func NewTestClient() (*Client, error) {
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = os.Getenv(accessSecretEnv)

	serverUrl, _ := url.Parse(server.URL)
	return NewClient(nil, serverUrl.String(), appToken, appSecret, accessToken, accessSecret)
}

// -----------Client tests-----------------------------------------------------

// Tests that a client is created fine
func TestClientCreation(t *testing.T) {
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = os.Getenv(accessSecretEnv)

	_, err := NewClient(nil, "http://resources/fake", appToken, appSecret, accessToken, accessSecret)

	if err != nil {
		t.Error("Error creating a client")
	}

}

// Tests the creation of the client with errors
func TestClientWrongParams(t *testing.T) {
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = ""

	_, err := NewClient(nil, "http://resources/fake", appToken, appSecret, accessToken, accessSecret)

	if err == nil {
		t.Error("Should be an error when creating the client")
	}
}

// Tests a creation default client
func TestDefaultClientCreation(t *testing.T) {
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = os.Getenv(accessSecretEnv)

	_, err := NewDefaultClient(appToken, appSecret, accessToken, accessSecret)

	if err != nil {
		t.Error("Error creating a default client")
	}

}

// Creates the client with error
func TestDefaultClientWrongParams(t *testing.T) {
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = ""

	_, err := NewDefaultClient(appToken, appSecret, accessToken, accessSecret)

	if err == nil {
		t.Error("Should be an error when creating the default client")
	}
}
