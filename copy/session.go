package copy

import (
	"errors"
	//"io"
	//"log"
	"net/http"
	"net/url"
	//"os"

	"github.com/garyburd/go-oauth/oauth"
)

const (
	ResourceURL = "https://api.copy.com/rest/"
	AuthURL     = "https://www.copy.com/applications/authorize"
)

type AccessToken struct {
	Token string
	Key   string
}

type AppToken struct {
	Token string
	Key   string
}

type Session struct {
	OauthClient oauth.Client
	TokenCreds  oauth.Credentials
}

func NewSession(appToken AppToken, accessToken AccessToken) Session {

	//Create app credentials
	appCreds := oauth.Credentials{
		Token:  appToken.Token,
		Secret: appToken.Key,
	}

	//Create the oauth client
	oauthClient := oauth.Client{
		TokenRequestURI: AuthURL,
		Credentials:     appCreds,
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
	req, err := http.NewRequest("GET", urlStr, nil)

	if err != nil {
		return nil, err
	}

	if req.URL.RawQuery != "" {
		return nil, errors.New("oauth: url must not contain a query string")
	}

	req.Header.Set("Authorization", s.OauthClient.AuthorizationHeader(&s.TokenCreds, "GET", req.URL, form))
	req.URL.RawQuery = form.Encode()

	return s.Do(req, nil)

}

func (s *Session) Post(urlStr string, form url.Values) (*http.Response, error) {
	return nil, nil
}

func (s *Session) Delete(urlStr string, form url.Values) (*http.Response, error) {
	return nil, nil
}

func (s *Session) Put(urlStr string, form url.Values) (*http.Response, error) {
	return nil, nil
}

func (s *Session) Do(request *http.Request, httpClient *http.Client) (*http.Response, error) {

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// Custom headers for Copy API, [IMPORTANT!!]
	customHeaders := map[string]string{
		"X-Api-Version": "1",
		"Accept":        "application/json",
	}

	for k, v := range customHeaders {
		request.Header.Add(k, v)
	}

	return httpClient.Do(request)

}
