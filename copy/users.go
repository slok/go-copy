package copy

import (
	"net/url"
)

// User represents the current user at Copy
type User struct {
	Id          string  `json:"id,omitempty"`
	FirstName   string  `json:"first_name,omitempty"`
	LastName    string  `json:"last_name,omitempty"`
	Developer   bool    `json:"developer,omitempty"`
	CreatedTime int     `json:"created_time,omitempty"`
	Email       string  `json:"email,omitempty"`
	Emails      []Email `json:"emails,omitempty"`
	Storage     Storage `json:"storage,omitempty"`
}

type Email struct {
	Primary   bool   `json:"primary,omitempty"`
	Confirmed bool   `json:"confirmed,omitempty"`
	Email     string `json:"Email,omitempty"`
	Gravatar  string `json:"gravatar,omitempty"`
}

type Storage struct {
	Used  int `json:"used,omitempty"`
	Quota int `json:"quota,omitempty"`
	Saved int `json:"saved,omitempty"`
}

type UserService struct {
	client *Client
}

const (
	endpointSuffix = "user"
)

func NewUserService(client *Client) *UserService {
	us := new(UserService)
	us.client = client
	return us
}

// Get fetches the authenticated user
//
//https://www.copy.com/developer/documentation#api-calls/profile
func (us *UserService) Get() (*User, error) {
	user := new(User)
	us.client.Do("GET", endpointSuffix, nil, user)
	return user, nil
}

// Updates the authenticated user
//
//https://www.copy.com/developer/documentation#api-calls/profile
func (us *UserService) Update(user *User) error {

	// Prepare the parameters to update (For now only frist and last name,
	// see copy docs)
	//
	// FIX: Don't craft by hand
	values := url.Values{
		"first_name": {user.FirstName},
		"last_name":  {user.LastName},
	}

	us.client.Do("PUT", endpointSuffix, values, user)
	return nil
}
