package copy

import ()

// User represents the current user at Copy
type User struct {
	Id          int     `json:"id,omitempty"`
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
	Used  int `json:"confirmed,omitempty"`
	Quota int `json:"quota,omitempty"`
	Saved int `json:"saved,omitempty"`
}

type UserService struct {
	client *Client
}

const (
	endpointSuffix = "user"
)

// Get fetches the authenticated user
//
//https://www.copy.com/developer/documentation#api-calls/profile
func (us *UserService) Get() (*User, error) {
	user := new(User)
	us.client.Do("GET", endpointSuffix, nil, user)
	return user, nil
}
