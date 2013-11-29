package copy

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var (
	userService *UserService
)

func setupUserService() {
	setup()
	userService = &UserService{client: client}
}

func tearDownUserService() {
	defer tearDown()
}

// Checks if the credentials for the integration tests are set in the env vars
func TestGetUser(t *testing.T) {
	setupUserService()
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
		t.Errorf("Users are not equal")
	}

}
