// Integration test
// This needs some env variables to run the tests
//      APP_TOKEN
//      APP_SECRET
//      ACCESS_TOKEN
//      ACCESS_SECRET
package copy

import (
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
)

func setup() {
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = os.Getenv(accessSecretEnv)
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

func TestGetRequest(t *testing.T) {

}
