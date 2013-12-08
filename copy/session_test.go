// Integration test
// This needs some env variables to run the tests
//      APP_TOKEN
//      APP_SECRET
//      ACCESS_TOKEN
//      ACCESS_SECRET

package copy

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	session *Session
)

func setupIntegration() {
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = os.Getenv(accessSecretEnv)

	session, _ = NewSession(
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

func tearDownIntegration() {

}

// Checks if the credentials for the integration tests are set in the env vars
func TestCredentialData(t *testing.T) {
	setupIntegration()
	defer tearDownIntegration()

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
	setupIntegration()
	defer tearDownIntegration()

	resp, err := session.Get(strings.Join([]string{defaultResourcesUrl, "user"}, "/"), nil, defaultHttpClient)

	if err != nil {
		t.Error("Expected no error in GET request")
	}

	if resp.StatusCode != 200 {
		t.Errorf("Response status error shouldn't be: %v", resp.StatusCode)
	}
}

// Check the GET request in an invalid copy resource with valid credentials
func TestGetRequestWrongResource(t *testing.T) {
	setupIntegration()
	defer tearDownIntegration()

	resp, _ := session.Get(strings.Join([]string{defaultResourcesUrl, "you shall not pass"}, "/"), nil, defaultHttpClient)

	if resp.StatusCode != 400 {
		t.Errorf("Response status error should be: %v", resp.StatusCode)
	}
}

// Check the GET request in a valid copy resource with wrong credentials
func TestGetRequestWrongCredentials(t *testing.T) {
	setupIntegration()
	defer tearDownIntegration()

	session.TokenCreds.Secret = "You shall not pass!"

	resp, _ := session.Get(strings.Join([]string{defaultResourcesUrl, "user"}, "/"), nil, defaultHttpClient)

	if resp.StatusCode != 400 {
		t.Errorf("Response status error should be: %v", resp.StatusCode)
	}
}

// Check the PUT request in a valid copy resource
func TestPutRequest(t *testing.T) {
	setupIntegration()
	defer tearDownIntegration()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	values := url.Values{
		"first_name": {fmt.Sprintf("TestName %d", r.Intn(100))},
		"last_name":  {fmt.Sprintf("TestSurname %d", r.Intn(100))},
	}

	resp, err := session.Put(strings.Join([]string{defaultResourcesUrl, "user"}, "/"), values, defaultHttpClient)

	defer resp.Body.Close()

	if err != nil {
		t.Error("Expected no error in POST request")
	}

	if resp.StatusCode != 200 {
		t.Errorf("Response status error shouldn't be: %v", resp.StatusCode)
	}

	// Get the response body to check the content
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)

	if err != nil {
		t.Errorf("Error reading the response body")
	}
	respBody := buf.String()

	if !strings.Contains(respBody, values["first_name"][0]) ||
		!strings.Contains(respBody, values["last_name"][0]) {
		t.Errorf("Not updated content with the REST API")
	}
}

// Check the PUT request in a valid copy resource
func TestDeleteRequest(t *testing.T) {
	setupIntegration()
	defer tearDownIntegration()

	// Put file with file service (we use client wrapper for convenience)
	client, err := NewDefaultClient(appToken, appSecret, accessToken, accessSecret)
	fs := NewFileService(client)
	err = fs.UploadFile("session_test.go", "session_test.go", true)
	if err != nil {
		t.Error("Could not prepare the DELETE integration test")
	}

	// Now test delete
	resp, err := session.Delete(strings.Join([]string{defaultResourcesUrl, "files", "session_test.go"}, "/"), nil, defaultHttpClient)
	resp.Body.Close()

	if err != nil {
		t.Error("Expected no error in Delete request")
	}

	if resp.StatusCode != 204 {
		t.Errorf("Response status error shouldn't be: %v", resp.StatusCode)
	}

	resp, err = session.Delete(strings.Join([]string{defaultResourcesUrl, "files", "doesntexists.go"}, "/"), nil, defaultHttpClient)

	if err == nil {
		t.Error("Expected error in Delete request")
	}
}
