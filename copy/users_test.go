package copy

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

var (
	userService *UserService
)

func setupUserService(t *testing.T) {
	setup(t)
	userService = NewUserService(client)
}

func tearDownUserService() {
	defer tearDown()
}

// Checks if the credentials for the integration tests are set in the env vars
func TestGetUser(t *testing.T) {
	setupUserService(t)
	defer tearDownUserService()

	mux.HandleFunc("/user",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w,
				`{
                  "id": "1381231",
                  "storage":{
                    "used": 9207643837,
                    "quota": 1100585369600,
                    "saved": 14557934927
                  },
                  "first_name": "Thomas",
                  "last_name": "Hunter",
                  "developer": true,
                  "created_time": 1358175510,
                  "email": "thomashunter@example.com",
                  "emails": [
                    {
                      "primary": true,
                      "confirmed": true,
                      "email": "thomashunter@example.com",
                      "gravatar": "eca957c6552e783627a0ced1035e1888"
                    },
                    {
                      "primary": false,
                      "confirmed": true,
                      "email": "thomashunter@example.net",
                      "gravatar": "c0e344ddcbabb383f94b1bd3486e55ba"
                    }
                  ]
                }
            `)
		},
	)

	user, _ := userService.Get()

	perfectUser := User{
		Id:          "1381231",
		FirstName:   "Thomas",
		LastName:    "Hunter",
		Developer:   true,
		CreatedTime: 1358175510,
		Email:       "thomashunter@example.com",
		Emails: []Email{
			Email{Primary: true,
				Confirmed: true,
				Email:     "thomashunter@example.com",
				Gravatar:  "eca957c6552e783627a0ced1035e1888"},
			Email{Primary: false,
				Confirmed: true,
				Email:     "thomashunter@example.net",
				Gravatar:  "c0e344ddcbabb383f94b1bd3486e55ba"},
		},
		Storage: Storage{
			Used:  9207643837,
			Quota: 1100585369600,
			Saved: 14557934927,
		},
	}

	// Are bouth content equal?
	if !reflect.DeepEqual(*user, perfectUser) {
		t.Errorf("objects are not equal")
	}

}

// Checks if the credentials for the integration tests are set in the env vars
func TestUpdateUser(t *testing.T) {
	setupUserService(t)
	defer tearDownUserService()

	//Our update data
	name := "Chuck"
	surname := "Norris"
	changeFlag := "Changed"

	perfectUser := User{
		FirstName: name,
		LastName:  surname,
	}

	mux.HandleFunc("/user",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")

			// Convert body to values
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			values, _ := url.ParseQuery(buf.String())

			// Get the data
			reqName := values["first_name"][0]
			reqSurname := values["last_name"][0]

			// Check if is the same data and if, seth the flag up
			if reqName != name {
				w.Header().Set("Status", "400 Bad Request")
			} else {
				reqName = reqName + changeFlag
			}

			if reqSurname != surname {
				w.Header().Set("Status", "400 Bad Request")
			} else {
				reqSurname = reqSurname + changeFlag
			}

			fmt.Fprintf(w,
				`{
                  "first_name": "%s",
                  "last_name": "%s"
              }`, reqName, reqSurname)
		},
	)

	userService.Update(&perfectUser)

	if perfectUser.FirstName != name+changeFlag || perfectUser.LastName != surname+changeFlag {
		t.Errorf("Could not update the user")
	}

}
