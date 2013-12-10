package copy

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var (
	fileService *FileService
)

func setupFileService(t *testing.T) {
	setup(t)
	fileService = NewFileService(client)
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

	// Test bad request
	server.Close()
	if _, err := fileService.GetTopLevelMeta(); err == nil {
		t.Errorf("No server up, should be an error")
	}

}

func TestGetMeta(t *testing.T) {
	setupFileService(t)
	defer tearDownFileService()

	mux.HandleFunc("/"+fmt.Sprintf(getMetaSuffix, "testing"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w,
				`{
               "id":"\/copy\/testing",
               "path":"\/testing",
               "name":"testing",
               "type":"dir",
               "size":null,
               "date_last_synced":1386150047,
               "modified_time":1386150047,
               "stub":false,
               "recipient_confirmed":false,
               "counts":[

               ],
               "mime_type":"",
               "link_name":null,
               "token":null,
               "creator_id":null,
               "permissions":null,
               "syncing":false,
               "public":false,
               "object_available":true,
               "links":[

               ],
               "url":"https:\/\/copy.com\/web\/users\/user-8129109\/copy\/testing",
               "thumb":null,
               "share":null,
               "children":[
                  {
                     "id":"\/copy\/testing\/random.txt",
                     "path":"\/testing\/random.txt",
                     "name":"random.txt",
                     "type":"file",
                     "size":1258291200,
                     "date_last_synced":1386151250,
                     "modified_time":1385993169,
                     "stub":true,
                     "recipient_confirmed":false,
                     "counts":[

                     ],
                     "mime_type":"text\/plain",
                     "link_name":null,
                     "token":null,
                     "creator_id":null,
                     "permissions":null,
                     "syncing":false,
                     "public":false,
                     "object_available":true,
                     "links":[

                     ],
                     "url":"https:\/\/copy.com\/web\/users\/user-8129109\/copy\/testing\/random.txt",
                     "revision":32,
                     "thumb":null,
                     "share":null,
                     "list_index":0
                  }
               ],
               "children_count":1
            }`)
		},
	)

	fileMeta, _ := fileService.GetMeta("testing")

	perfectFileMeta := Meta{
		Id:                 "/copy/testing",
		Path:               "/testing",
		Name:               "testing",
		Type:               "dir",
		DateLastSynced:     1386150047,
		ModifiedTime:       1386150047,
		Stub:               false,
		RecipientConfirmed: false,
		Syncing:            false,
		Public:             false,
		ObjectAvailable:    true,
		Url:                "https://copy.com/web/users/user-8129109/copy/testing",
		Links:              []Link{}, // for Deep equal nil and empty slice aren't the same
		Children: []Meta{
			Meta{
				Id:                 "/copy/testing/random.txt",
				Path:               "/testing/random.txt",
				Name:               "random.txt",
				Type:               "file",
				Size:               1258291200,
				DateLastSynced:     1386151250,
				ModifiedTime:       1385993169,
				Stub:               true,
				RecipientConfirmed: false,
				MimeType:           "text/plain",
				Syncing:            false,
				Public:             false,
				ObjectAvailable:    true,
				Url:                "https://copy.com/web/users/user-8129109/copy/testing/random.txt",
				Revision:           32,
				ListIndex:          0,
				Links:              []Link{}, // for Deep equal nil and empty slice aren't the same

			},
		},
		ChildrenCount: 1,
	}

	// Are bouth content equal?
	if !reflect.DeepEqual(*fileMeta, perfectFileMeta) {
		t.Errorf("Metas are not equal")
	}

	// Test bad request
	server.Close()
	if _, err := fileService.GetTopLevelMeta(); err == nil {
		t.Errorf("No server up, should be an error")
	}

}

func TestListRevisionsMeta(t *testing.T) {
	setupFileService(t)
	defer tearDownFileService()

	mux.HandleFunc("/"+fmt.Sprintf(listRevisionsSuffix, "Big API hanges/API-Changes.md"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w,
				`{
                      "id": "/copy/Big%20API%20Changes/API-Changes.md/@activity",
                      "path": "/Big API Changes/API-Changes.md",
                      "name": "Activity",
                      "token": null,
                      "permissions": null,
                      "syncing": false,
                      "public": false,
                      "type": "file",
                      "size": 12670,
                      "date_last_synced": 1365543105,
                      "stub": false,
                      "recipient_confirmed": false,
                      "url": "https://copy.com/web/Big%20API%20Changes/API-Changes.md",
                      "revision_id": "5000",
                      "thumb": null,
                      "share": null,
                      "counts": [
                      ],
                      "links": [
                      ],
                      "revisions": [
                        {
                          "revision_id": "5000",
                          "modified_time": "1365543105",
                          "size": 12670,
                          "latest": true,
                          "conflict": 4324,
                          "id": "/copy/Big%20API%20Changes/API-Changes.md/@activity/@time:1365543105",
                          "type": "revision",
                          "creator": {
                            "user_id": "1381231",
                            "created_time": 1358175510,
                            "email": "thomashunter@example.com",
                            "first_name": "Thomas",
                            "last_name": "Hunter",
                            "confirmed": true
                          }
                        },
                        {
                          "revision_id": "4900",
                          "modified_time": "1365542000",
                          "size": 12661,
                          "latest": false,
                          "conflict": 4324,
                          "id": "/copy/Big%20API%20Changes/API-Changes.md/@activity/@time:1365542000",
                          "type": "revision",
                          "creator": {
                            "user_id": "1381231",
                            "created_time": 1358175510,
                            "email": "thomashunter@example.com",
                            "first_name": "Thomas",
                            "last_name": "Hunter",
                            "confirmed": true
                          }
                        },
                        {
                          "revision_id": "4800",
                          "modified_time": "1365543073",
                          "size": 12658,
                          "latest": false,
                          "conflict": 4324,
                          "id": "/copy/Big%20API%20Changes/API-Changes.md/@activity/@time:1365543073",
                          "type": "revision",
                          "creator": {
                            "user_id": "1381231",
                            "created_time": 1358175510,
                            "email": "thomashunter@example.com",
                            "first_name": "Thomas",
                            "last_name": "Hunter",
                            "confirmed": true
                          }
                        }
                      ]
                    }
                `)
		},
	)

	revisions, _ := fileService.ListRevisionsMeta("Big API hanges/API-Changes.md")
	perfectRevisions := []Revision{
		Revision{
			RevisionId:   "5000",
			ModifiedTime: "1365543105",
			Size:         12670,
			Latest:       true,
			Conflict:     4324,
			Id:           "/copy/Big%20API%20Changes/API-Changes.md/@activity/@time:1365543105",
			Type:         "revision",
			Creator: Creator{
				UserId:      "1381231",
				CreatedTime: 1358175510,
				Email:       "thomashunter@example.com",
				FirstName:   "Thomas",
				LastName:    "Hunter",
				Confirmed:   true,
			},
		},
		Revision{
			RevisionId:   "4900",
			ModifiedTime: "1365542000",
			Size:         12661,
			Latest:       false,
			Conflict:     4324,
			Id:           "/copy/Big%20API%20Changes/API-Changes.md/@activity/@time:1365542000",
			Type:         "revision",
			Creator: Creator{
				UserId:      "1381231",
				CreatedTime: 1358175510,
				Email:       "thomashunter@example.com",
				FirstName:   "Thomas",
				LastName:    "Hunter",
				Confirmed:   true,
			},
		},
		Revision{
			RevisionId:   "4800",
			ModifiedTime: "1365543073",
			Size:         12658,
			Latest:       false,
			Conflict:     4324,
			Id:           "/copy/Big%20API%20Changes/API-Changes.md/@activity/@time:1365543073",
			Type:         "revision",
			Creator: Creator{
				UserId:      "1381231",
				CreatedTime: 1358175510,
				Email:       "thomashunter@example.com",
				FirstName:   "Thomas",
				LastName:    "Hunter",
				Confirmed:   true,
			},
		},
	}

	// Are bouth content equal?
	if !reflect.DeepEqual(revisions, perfectRevisions) {
		t.Errorf("Metas are not equal")
	}

	// Test bad request
	server.Close()
	if _, err := fileService.ListRevisionsMeta("Big API hanges/API-Changes.md"); err == nil {
		t.Errorf("No server up, should be an error")
	}
}

// Checks json decoding for the meta object
func TestGetFile(t *testing.T) {
	setupFileService(t)
	defer tearDownFileService()

	filename := "client_test.go"
	// Read the file to test
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		t.Error(err.Error())
	}

	mux.HandleFunc(strings.Join([]string{"", filesTopLevelSuffix, filename}, "/"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			w.Write(file)
		},
	)

	fileReader, _ := fileService.GetFile(filename)
	defer fileReader.Close()

	file2, err := ioutil.ReadAll(fileReader)

	if err != nil {
		t.Error(err.Error())
	}

	if !bytes.Equal(file, file2) {
		t.Errorf("contents are not equal")
	}

	// Test bad request
	server.Close()
	if _, err := fileService.GetFile(filename); err == nil {
		t.Errorf("No server up, should be an error")
	}
}

func TestFileUpload(t *testing.T) {
	setupFileService(t)
	defer tearDownFileService()

	filePath := "files_test.go"
	upPath := "tests/uploads"

	// Read the file to test
	origFile, err := ioutil.ReadFile(filePath)

	if err != nil {
		t.Error(err.Error())
	}

	resPath := strings.Join([]string{"", filesTopLevelSuffix, upPath}, "/")

	mux.HandleFunc(resPath,
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			// Check that upload is ok
			r.ParseMultipartForm(100000)
			form := r.MultipartForm

			files, _ := form.File["file"]
			file, _ := files[0].Open()
			defer file.Close()

			buf := new(bytes.Buffer)
			io.Copy(buf, file)

			if !bytes.Equal(origFile, buf.Bytes()) {
				t.Errorf("contents are not equal")
			}
		},
	)

	err = fileService.UploadFile(filePath, strings.Join([]string{upPath, filePath}, "/"), true)

	if err != nil {
		t.Error(err.Error())
	}

	// Test bad request
	server.Close()
	if err := fileService.UploadFile(filePath, strings.Join([]string{upPath, filePath}, "/"), true); err == nil {
		t.Errorf("No server up, should be an error")
	}
}

func TestRenameFile(t *testing.T) {

	setupFileService(t)
	defer tearDownFileService()

	filePath := "test/test2"
	newName := "test2.2"
	overwrite := true

	regex := "/" + filesTopLevelSuffix + `/(.+)\?name=(.+)&overwrite=(.*)`
	re, _ := regexp.Compile(regex)

	mux.HandleFunc(strings.Join([]string{"", filesTopLevelSuffix, filePath}, "/"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")

			matches := re.FindAllStringSubmatch(r.URL.String(), -1)
			path := matches[0][1]
			name := matches[0][2]

			ow := false

			if matches[0][3] == "true" {
				ow = true
			}

			if filePath != path || newName != name || overwrite != ow {
				t.Errorf("Wrong params in URL")
			}

		},
	)

	if err := fileService.RenameFile(filePath, newName, overwrite); err != nil {
		t.Errorf("Shouldn't be an error")
	}

	// Test bad request
	server.Close()
	if err := fileService.RenameFile(filePath, newName, overwrite); err == nil {
		t.Errorf("No server up, should be an error")
	}
}

func TestMoveFile(t *testing.T) {

	setupFileService(t)
	defer tearDownFileService()

	filePath := "test/test2.txt"
	newPath := "test3/test2.txt"
	overwrite := true

	regex := "/" + filesTopLevelSuffix + `/(.+)\?path=(.+)&overwrite=(.+)`
	re, _ := regexp.Compile(regex)
	mux.HandleFunc(strings.Join([]string{"", filesTopLevelSuffix, filePath}, "/"),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")
			matches := re.FindAllStringSubmatch(r.URL.String(), -1)
			path := matches[0][1]
			movePath := matches[0][2]

			ow := false

			if matches[0][3] == "true" {
				ow = true
			}

			if filePath != path || newPath != movePath || overwrite != ow {
				t.Errorf("Wrong params in URL")
			}

		},
	)

	if err := fileService.MoveFile(filePath, newPath, overwrite); err != nil {
		t.Errorf("Shouldn't be an error")
	}

	// Test bad request
	server.Close()
	if err := fileService.MoveFile(filePath, newPath, overwrite); err == nil {
		t.Errorf("No server up, should be an error")
	}
}

func TestDeleteFile(t *testing.T) {

	setupFileService(t)
	defer tearDownFileService()

	filePath := "test/test2.txt"

	resPath := strings.Join([]string{"", filesTopLevelSuffix, filePath}, "/")
	regex := "/" + filesTopLevelSuffix + `/(.+)`
	re, _ := regexp.Compile(regex)

	mux.HandleFunc(resPath,
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "DELETE")

			matches := re.FindAllStringSubmatch(r.URL.String(), -1)
			path := matches[0][1]

			if path != filePath {
				t.Errorf("Wrong params in URL")
			}

		},
	)

	if err := fileService.DeleteFile(filePath); err != nil {
		t.Errorf("Shouldn't be an error")
	}

	// Test bad request
	server.Close()
	if err := fileService.DeleteFile(filePath); err == nil {
		t.Errorf("No server up, should be an error")
	}
}

func TestCreateDirFile(t *testing.T) {

	setupFileService(t)
	defer tearDownFileService()

	filePath := "dir1/dir2"
	overwrite := true

	resPath := strings.Join([]string{"", filesTopLevelSuffix, filePath}, "/")
	regex := "/" + filesTopLevelSuffix + `/(.+)\?overwrite=(.+)`
	re, _ := regexp.Compile(regex)

	mux.HandleFunc(resPath,
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			matches := re.FindAllStringSubmatch(r.URL.String(), -1)
			path := matches[0][1]

			ow := false

			if matches[0][2] == "true" {
				ow = true
			}

			if filePath != path || overwrite != ow {
				t.Errorf("Wrong params in URL")
			}

		},
	)

	if err := fileService.CreateDirectory(filePath, overwrite); err != nil {
		t.Errorf("Shouldn't be an error")
	}

	// Test bad request
	server.Close()
	if err := fileService.CreateDirectory(filePath, overwrite); err == nil {
		t.Errorf("No server up, should be an error")
	}

}

func TestGetRevisionMeta(t *testing.T) {
	setupFileService(t)
	defer tearDownFileService()

	mux.HandleFunc("/"+fmt.Sprintf(revisionSuffix, "Big API hanges/API-Changes.md", 1365532651),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w,
				`{
                    "id": "/copy/Big%20API%20Changes/API-Changes.md/@activity/@time:1365532651",
                    "path": "/Big API Changes/API-Changes.md",
                    "name": "API-Changes.md",
                    "token": null,
                    "permissions": null,
                    "syncing": false,
                    "public": false,
                    "type": "file",
                    "size": 12666,
                    "date_last_synced": 1365532651,
                    "stub": false,
                    "recipient_confirmed": false,
                    "url": "https://copy.com/web/Big%20API%20Changes/API-Changes.md?revision=4898",
                    "revision_id": 4898,
                    "thumb": null,
                    "share": null,
                    "counts": [
                    ],
                    "links": [
                      ]
                }
                `)
		},
	)

	revision, _ := fileService.GetRevisionMeta("Big API hanges/API-Changes.md", 1365532651)
	perfectRevision := Meta{

		Id:                 "/copy/Big%20API%20Changes/API-Changes.md/@activity/@time:1365532651",
		Path:               "/Big API Changes/API-Changes.md",
		Name:               "API-Changes.md",
		Syncing:            false,
		Public:             false,
		Type:               "file",
		Size:               12666,
		DateLastSynced:     1365532651,
		Stub:               false,
		RecipientConfirmed: false,
		Url:                "https://copy.com/web/Big%20API%20Changes/API-Changes.md?revision=4898",
		RevisionId:         4898,
		Counts:             Count{},
		Links:              []Link{},
	}

	// Are bouth content equal?
	if !reflect.DeepEqual(*revision, perfectRevision) {
		t.Errorf("Metas are not equal")
	}

	// Test bad request
	server.Close()
	if _, err := fileService.ListRevisionsMeta("Big API hanges/API-Changes.md"); err == nil {
		t.Errorf("No server up, should be an error")
	}
}
