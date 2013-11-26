package copy

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/garyburd/go-oauth/oauth"
)

const (
	ResourceURL = "https://api.copy.com/rest/"
	AuthURL     = "https://www.copy.com/applications/authorize"
)

type AccessToken struct {
	Key   string
	Token string
}

type Session struct {
	OauthClient oauth.Client
	TokenCreds  oauth.Credentials
}

func NewSession(accessToken AccessToken) Session {

	//Create the oauth client
	oauthClient := oauth.Client{
		TokenRequestURI: AuthURL,
	}

	tokenCred := oauth.Credentials{
		Token:  accessToken.Token,
		Secret: accessToken.Key,
	}

	//Return Session with the Oauth client created
	return Session{
		OauthClient: oauthClient,
		TokenCreds:  tokenCred,
	}

}

func (s *Session) Get(urlStr string, form url.Values) (*http.Response, error) {

}

func (s *Session) Post(urlStr string, form url.Values) (*http.Response, error) {

}

func (s *Session) Delete(urlStr string, form url.Values) (*http.Response, error) {

}

func (s *Session) Put(urlStr string, form url.Values) (*http.Response, error) {

}

func (s *Session) Do(request *http.Request) (*http.Response, error) {

	// Custom headers for Copy API
	customHeades := map[string]string{
		"X-Api-Version": "1",
		"Accept":        "application/json",
	}

}
