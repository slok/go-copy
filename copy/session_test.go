// Integration test
// This needs some env variables to run the tests
//      APP_TOKEN
//      APP_SECRET
//      ACCESS_TOKEN
//      ACCESS_SECRET
package copy

import (
	"io"
	"log"
	"net/url"
	"os"
	"testing"
)

const (
	appTokenEnv     = "APP_TOKEN"
	appSecretEnv    = "APP_SECRET"
	accessTokenEnv  = "ACCESS_TOKEN"
	accessSecretEnv = "ACCESS_SECRET"
)

var (
	appToken     string
	appSecret    string
	accessToken  string
	accessSecret string
	session      Session
)

func setup() {
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = os.Getenv(accessSecretEnv)

	session = NewSession(
		AppToken{
			Token: appToken,
			Key:   appSecret,
		},
		AccessToken{
			Token: accessToken,
			Key:   accessSecret,
		},
	)
}

func tearDown() {

}

// Checks if the credentials for the integration tests are set in the env vars
func TestCredentialData(t *testing.T) {
	setup()
	defer tearDown()

	if appToken == "" {
		t.Error("Expected", appTokenEnv, "env var")
	}

	if appSecret == "" {
		t.Error("Expected", appSecretEnv, "env var")
	}

	if accessToken == "" {
		t.Error("Expected", accessTokenEnv, "env var")
	}

	if accessSecret == "" {
		t.Error("Expected", accessSecretEnv, "env var")
	}
}

// Check the GET request in a valid copy resource
func TestGetRequest(t *testing.T) {
	setup()
	defer tearDown()

	resp, err := session.Get("https://api.copy.com/rest/user", nil)

	if err != nil {
		t.Error("Expected no error in GET request")
	}

	if resp.StatusCode != 200 {
		t.Errorf("Response status error shouldn't be: %v", resp.StatusCode)
	}
}

// Check the GET request in an invalid copy resource with valid credentials
func TestGetRequestWrongResource(t *testing.T) {
	setup()
	defer tearDown()

	resp, _ := session.Get("https://api.copy.com/rest/userfail", nil)

	if resp.StatusCode != 400 {
		t.Errorf("Response status error should be: %v", resp.StatusCode)
	}
}

// Check the GET request in a valid copy resource with wrong credentials
func TestGetRequestWrongCredentials(t *testing.T) {
	setup()
	defer tearDown()

	session.TokenCreds.Secret = "You shall not pass!"

	resp, _ := session.Get("https://api.copy.com/rest/user", nil)

	if resp.StatusCode != 400 {
		t.Errorf("Response status error should be: %v", resp.StatusCode)
	}
}

// Check the PUT request in a valid copy resource
func TestPutRequest(t *testing.T) {
	setup()
	defer tearDown()

	values := url.Values{}
	values.Add("first_name", "Slok2")
	values.Add("last_name", "Wooloo")

	resp, err := session.Put("https://api.copy.com/rest/user", values)

	defer resp.Body.Close()
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		t.Error("Expected no error in POST request")
	}

	if resp.StatusCode != 200 {
		t.Errorf("Response status error shouldn't be: %v", resp.StatusCode)
	}
}
