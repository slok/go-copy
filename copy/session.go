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
)

type AccessToken struct {
	Key   string
	Token string
}

type Session struct {
	AccessToken AccessToken
}

func NewSession(accessToken AccessToken) Session {
	//Return Session
	return Session{
		AccessToken{
			Key:   accessToken.Key,
			Token: accessToken.Token,
		},
	}

}

func (s *Session) Get() {

}

func (s *Session) Post() {

}

func (s *Session) Delete() {

}

func (s *Session) Put() {

}

func (s *Session) Do() {

}
