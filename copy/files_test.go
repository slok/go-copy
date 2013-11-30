package copy

import (
	//"bytes"
	"fmt"
	"net/http"
	//"net/url"
	//"os"
	"reflect"
	"testing"
)

var (
	fileService *FileService
)

func setupFileService(t *testing.T) {
	setup(t)
	fileService = &FileService{client: client}
}

func tearDownFileService() {
	defer tearDown()
}

// Checks if the credentials for the integration tests are set in the env vars
func TestGetTopLevelMeta(t *testing.T) {
	setupFileService(t)
	defer tearDownFileService()
	mux.HandleFunc("/meta",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w,
				`{
               "id":"\/",
               "path":"\/",
               "name":"Copy",
               "type":"root",
               "stub":false,
               "children":[
                  {
                     "name":"Personal Files",
                     "type":"copy",
                     "id":"\/copy",
                     "path":"\/",
                     "stub":true,
                     "counts":{
                        "new":0,
                        "viewed":0,
                        "hidden":0
                     }
                  }
               ],
               "children_count":1
            }`)
		},
	)

	fileMeta, _ := fileService.GetTopLevelMeta()

	perfectFileMeta := Meta{
		Id:   "/",
		Path: "/",
		Name: "Copy",
		Type: "root",
		Stub: false,
		Children: []Meta{
			Meta{
				Id:   "/copy",
				Path: "/",
				Name: "Personal Files",
				Type: "copy",
				Stub: true,
				Counts: Count{
					New:    0,
					Viewed: 0,
					Hidden: 0,
				},
			},
		},
		ChildrenCount: 1,
	}

	// Are bouth content equal?
	if !reflect.DeepEqual(*fileMeta, perfectFileMeta) {
		t.Errorf("Metas are not equal")
	}

	/*
		//Prepare the neccesary data
		appToken := os.Getenv("APP_TOKEN")
		appSecret := os.Getenv("APP_SECRET")
		accessToken := os.Getenv("ACCESS_TOKEN")
		accessSecret := os.Getenv("ACCESS_SECRET")

		// Create the client
		client, err := NewDefaultClient(appToken, appSecret, accessToken, accessSecret)
		if err != nil {
			fmt.Fprint(os.Stderr, "Could not create the client, review the auth params")
			os.Exit(-1)
		}

		//Create the service (in this case for a user)
		fileService = NewFileService(client)

		//Play with the lib :)
		fileMeta, err := fileService.GetTopLevelMeta()
		if err != nil {
			fmt.Fprint(os.Stderr, "Could not retrieve the user")
			os.Exit(-1)
		}
	*/
	//Print the object with reflection (used for debugging)
	/*val := reflect.ValueOf(fileMeta).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		fmt.Printf("%s\t: %v\n", typeField.Name, valueField.Interface())
	}*/
}
