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

// Checks json decoding for the meta object
func TestJsonMetaDecoding(t *testing.T) {
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
               "children_count":1,
               "link_name":"link test",
               "token":"32234dsad",
               "permissions":"all",
               "public":true,
               "size":3123123,
               "date_last_synced":32131232,
               "share":true,
               "recipient_confirmed":true,
               "object_available":true,
               "links": [
                   {
                        "id":"link1",
                        "public":true,
                        "expires":true,
                        "expired":true,
                        "url":"dsafdsfdsaxfwf",
                        "url_short":"dsadsad",
                        "recipients": [
                            {
                                "contact_Type":"gfgdfd",
                                "contact_id":"fgffsd",
                                "contact_source":"htgdffvdb",
                                "user_id":"3343",
                                "first_name":"ffgfgf",
                                "last_name":"grfesa",
                                "email":"fsdfdsfds",
                                "permissions":"all",
                                "emails": [
                                     {
                                            "confirmed":true,
                                            "primary":true,
                                            "email":"thomashunter@example.com",
                                            "gravatar":"eca957c6552e783627a0ced1035e1888"
                                    }
                                ]
                            }
                        ],
                        "creator_id":"htgdffsdd",
                        "confirmation_required": true
                    }
               ],
               "revisions": [
                    {
                        "revision_id":"231312",
                        "modified_time":"32324",
                        "size":31232,
                        "latest":true,
                        "conflict":4324,
                        "id":"dsdsd",
                        "type":"sdsad",
                        "creator":{
                            "user_id":"44342",
                            "created_time":323423,
                            "email":"fdfdsf@dsadsa.com",
                            "first_name":"sadasd",
                            "last_name":"sdsadsafds",
                            "confirmed":true
                        }
                    }
                ],
                "url":"dasdsafdasddfdf",
                "revision_id":31312,
                "thumb":"test thumb",
                "thumb_original_dimensions":{
                    "width":32432,
                    "height":53543
                }
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
		ChildrenCount:      1,
		LinkName:           "link test",
		Token:              "32234dsad",
		Permissions:        "all",
		Public:             true,
		Size:               3123123,
		DateLastSynced:     32131232,
		Share:              true,
		RecipientConfirmed: true,
		ObjectAvailable:    true,
		Links: []Link{
			Link{
				Id:       "link1",
				Public:   true,
				Expires:  true,
				Expired:  true,
				Url:      "dsafdsfdsaxfwf",
				UrlShort: "dsadsad",
				Recipients: []Recipient{
					Recipient{
						ContactType:   "gfgdfd",
						ContactId:     "fgffsd",
						ContactSource: "htgdffvdb",
						UserId:        "3343",
						FirstName:     "ffgfgf",
						LastName:      "grfesa",
						Email:         "fsdfdsfds",
						Permissions:   "all",
						Emails: []Email{
							Email{
								Confirmed: true,
								Primary:   true,
								Email:     "thomashunter@example.com",
								Gravatar:  "eca957c6552e783627a0ced1035e1888",
							},
						},
					},
				},
				CreatorId:            "htgdffsdd",
				ConfirmationRequired: true,
			},
		},
		Revisions: []Revision{
			Revision{
				RevisionId:   "231312",
				ModifiedTime: "32324",
				Size:         31232,
				Latest:       true,
				Conflict:     4324,
				Id:           "dsdsd",
				Type:         "sdsad",
				Creator: Creator{
					UserId:      "44342",
					CreatedTime: 323423,
					Email:       "fdfdsf@dsadsa.com",
					FirstName:   "sadasd",
					LastName:    "sdsadsafds",
					Confirmed:   true,
				},
			},
		},
		Url:        "dasdsafdasddfdf",
		RevisionId: 31312,
		Thumb:      "test thumb",
		ThumbOriginalDimensions: ThumbOriginalDimensions{
			Width:  32432,
			Height: 53543,
		},
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
