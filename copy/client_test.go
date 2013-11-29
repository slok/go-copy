package copy

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

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
	serverUrl, _ := url.Parse(server.URL)

	var err error
	client, err = NewClient(nil, serverUrl.String())
	if err != nil {
		t.Error(err.Error())
	}
}

func tearDown() {
	server.Close()
}

// From go-github (https://github.com/google/go-github)
func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}
