package copy

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"testing"
)

// Global testing vars
var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

// Setups the mock server for all the app tests
func setup(t *testing.T) {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	var err error
	client, err = NewTestClient()
	if err != nil {
		t.Error(err.Error())
	}
}

// ---------Global testing utils-----------------------------------------------

// Cleans all the setup of the test
func tearDown() {
	server.Close()
}

// From go-github (https://github.com/google/go-github)
func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

// Creates a new Client for testing
func NewTestClient() (*Client, error) {
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = os.Getenv(accessSecretEnv)

	serverUrl, _ := url.Parse(server.URL)
	return NewClient(nil, serverUrl.String(), appToken, appSecret, accessToken, accessSecret)
}

// -----------Client tests-----------------------------------------------------

// Tests that a client is created fine
func TestClientCreation(t *testing.T) {
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = os.Getenv(accessSecretEnv)

	_, err := NewClient(http.DefaultClient, "http://resources/fake", appToken, appSecret, accessToken, accessSecret)

	if err != nil {
		t.Error("Error creating a client")
	}

}

// Tests the creation of the client with errors
func TestClientWrongParams(t *testing.T) {
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = ""

	_, err := NewClient(http.DefaultClient, "http://resources/fake", appToken, appSecret, accessToken, accessSecret)

	if err == nil {
		t.Error("Should be an error when creating the client")
	}
}

// Tests a creation default client
func TestDefaultClientCreation(t *testing.T) {
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = os.Getenv(accessSecretEnv)

	_, err := NewDefaultClient(appToken, appSecret, accessToken, accessSecret)

	if err != nil {
		t.Error("Error creating a default client")
	}

}

// Creates the client with error
func TestDefaultClientWrongParams(t *testing.T) {
	appToken = os.Getenv(appTokenEnv)
	appSecret = os.Getenv(appSecretEnv)
	accessToken = os.Getenv(accessTokenEnv)
	accessSecret = ""

	_, err := NewDefaultClient(appToken, appSecret, accessToken, accessSecret)

	if err == nil {
		t.Error("Should be an error when creating the default client")
	}
}

// Client request methods tests------------------------------------------------
type testObject struct {
	Field1 string        `json:"field1,omitempty"`
	Field2 bool          `json:"field2,omitempty"`
	Field3 int           `json:"field3,omitempty"`
	Field4 float32       `json:"field4,omitempty"`
	Field5 []testObject2 `json:"field5,omitempty"`
}

type testObject2 struct {
	Field1b string `json:"field1b,omitempty"`
	Field2b []int  `json:"field2b,omitempty"`
}

// Our test object
var perfectObject = testObject{
	Field1: "testfield1",
	Field2: true,
	Field3: 26,
	Field4: 2.6,
	Field5: []testObject2{
		testObject2{
			Field1b: "testfield1b",
			Field2b: []int{
				0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			},
		},
		testObject2{
			Field1b: "testfield1b2",
			Field2b: []int{
				10, 11, 12, 13, 14, 15,
			},
		},
	},
}

var jsonObject = `{
	"field1": "testfield1",
	"field2": true,
	"field3": 26,
	"field4": 2.6,
	"field5": [
		{
			"field1b": "testfield1b",
			"field2b": [0,1,2,3,4,5,6,7,8,9]
		},
		{
			"field1b": "testfield1b2",
			"field2b": [10,11,12,13,14,15]
		}
	]
}`

// Check that the request decoding is ok with GET
func TestDoRequestDecodingGET(t *testing.T) {
	// Prepare the mock server
	setup(t)
	defer tearDown()

	mux.HandleFunc("/do-request-decoding",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, jsonObject)
		},
	)

	obj := new(testObject)
	client.DoRequestDecoding("GET", "do-request-decoding", nil, obj)

	// Are both content equal?
	if !reflect.DeepEqual(*obj, perfectObject) {
		t.Errorf("Objects are not equal")
	}
}

// Check that the request decoding is ok with PUT
func TestDoRequestDecodingPUT(t *testing.T) {
	// Prepare the mock server
	setup(t)
	defer tearDown()

	mux.HandleFunc("/do-request-decoding",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")

			// Convert body to values
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			values, _ := url.ParseQuery(buf.String())

			// Get the data
			f1 := values["field1"][0]
			f3, _ := strconv.Atoi(values["field3"][0])

			fmt.Fprintf(w,
				`{"field1": "%v",
				"field2": true,
				"field3": %d,
				"field4": 2.6,
				"field5": [
					{
						"field1b": "testfield1b",
						"field2b": [0,1,2,3,4,5,6,7,8,9]
					},
					{
						"field1b": "testfield1b2",
						"field2b": [10,11,12,13,14,15]
					}
				]
			}`, f1, f3)
		},
	)

	values := url.Values{
		"field1": []string{"value1"},
		"field3": []string{"234"},
	}

	obj := new(testObject)
	client.DoRequestDecoding("PUT", "do-request-decoding", values, obj)

	// Are both content equal?
	perfectObject.Field1 = values["field1"][0]
	perfectObject.Field3, _ = strconv.Atoi(values["field3"][0])

	if !reflect.DeepEqual(*obj, perfectObject) {
		t.Errorf("Objects are not equal")
	}
}

// Check wrong request
func TestDoRequestDecodingWrongMethod(t *testing.T) {
	// Prepare the mock server
	setup(t)
	defer tearDown()

	obj := new(testObject)
	_, err := client.DoRequestDecoding("FAKE", "do-request-decoding", nil, obj)

	if err == nil {
		t.Errorf("FAKE method should raise an error")
	}

}
