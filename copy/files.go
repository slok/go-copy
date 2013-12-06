package copy

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// File Meta data representation
type Meta struct {
	Id                      string                  `json:"id,omitempty"`
	Path                    string                  `json:"path,omitempty"`
	Name                    string                  `json:"name,omitempty"`
	LinkName                string                  `json:"link_name,omitempty"`
	Token                   string                  `json:"token,omitempty"`
	Permissions             string                  `json:"permissions,omitempty"`
	Public                  bool                    `json:"public,omitempty"`
	Type                    string                  `json:"type,omitempty"`
	Size                    int                     `json:"size,omitempty"`
	DateLastSynced          int                     `json:"date_last_synced,omitempty"`
	ModifiedTime            int                     `json:"modified_time,omitempty"`
	Stub                    bool                    `json:"stub,omitempty"`
	Share                   bool                    `json:"share,omitempty"`
	Children                []Meta                  `json:"children,omitempty"` // Inception :D
	Counts                  Count                   `json:"counts,omitempty"`   // Array? (sometimes? ask copy.com)
	RecipientConfirmed      bool                    `json:"recipient_confirmed",omitempty"`
	MimeType                string                  `json:"mime_type",omitempty"`
	Syncing                 bool                    `json:"syncing",omitempty"`
	ObjectAvailable         bool                    `json:"object_available,omitempty"`
	Links                   []Link                  `json:"links,omitempty"`
	Revisions               []Revision              `json:"revisions,omitempty"`
	Url                     string                  `json:"url,omitempty"`
	RevisionId              int                     `json:"revision_id,omitempty"`
	Thumb                   string                  `json:"thumb,omitempty"`
	ThumbOriginalDimensions ThumbOriginalDimensions `json:"thumb_original_dimensions,omitempty"`
	ChildrenCount           int                     `json:"children_count",omitempty"`
	Revision                int                     `json:"revision",omitempty"`
	ListIndex               int                     `json:"list_index",omitempty"`
}

type Count struct {
	New    int `json:"new,omitempty"`
	Viewed int `json:"viewed,omitempty"`
	Hidden int `json:"hidden,omitempty"`
}

type Link struct {
	Id                   string      `json:"id,omitempty"`
	Public               bool        `json:"public,omitempty"`
	Expires              bool        `json:"expires,omitempty"`
	Expired              bool        `json:"expired,omitempty"`
	Url                  string      `json:"url,omitempty"`
	UrlShort             string      `json:"url_short,omitempty"`
	Recipients           []Recipient `json:"recipients,omitempty"`
	CreatorId            string      `json:"creator_id,omitempty"`
	ConfirmationRequired bool        `json:"confirmation_required,omitempty"`
}
type Recipient struct {
	ContactType   string  `json:"contact_type,omitempty"`
	ContactId     string  `json:"contact_id,omitempty"`
	ContactSource string  `json:"contact_source,omitempty"`
	UserId        string  `json:"user_id,omitempty"`
	FirstName     string  `json:"first_name,omitempty"`
	LastName      string  `json:"last_name,omitempty"`
	Email         string  `json:"email,omitempty"`
	Permissions   string  `json:"permissions,omitempty"`
	Emails        []Email `json:"emails,omitempty"` // In users.go
}

type ThumbOriginalDimensions struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"Height,omitempty"`
}

type Revision struct {
	RevisionId   string  `json:"revision_id,omitempty"`
	ModifiedTime string  `json:"modified_time,omitempty"`
	Size         int     `json:"size,omitempty"`
	Latest       bool    `json:"latest,omitempty"`
	Conflict     int     `json:"conflict,omitempty"`
	Id           string  `json:"id,omitempty"`
	Type         string  `json:"type,omitempty"`
	Creator      Creator `json:"creator,omitempty"`
}

type Creator struct {
	UserId      string `json:"user_id,omitempty"`
	CreatedTime int    `json:"created_time,omitempty"`
	Email       string `json:"email,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Confirmed   bool   `json:"confirmed,omitempty"`
}

type FileService struct {
	client *Client
}

var (
	//Options
	overwriteOption = "overwrite=%t"
	nameOption      = "name=%v"
	pathOption      = "path=%v"
	sizeOption      = "size=%d"

	// Meta paths
	metaTopLevelSuffix  = "meta"                                                        // http.../meta/PATH
	firstLevelSuffix    = strings.Join([]string{metaTopLevelSuffix, "copy"}, "/")       // http.../meta/copy
	getMetaSuffix       = strings.Join([]string{firstLevelSuffix, "%v"}, "/")           // http.../meta/copy/PATH
	listRevisionsSuffix = strings.Join([]string{firstLevelSuffix, "%v/@activity"}, "/") // http.../meta/copy/PATH/@activity
	revisionSuffix      = strings.Join([]string{listRevisionsSuffix, "@time:%d"}, "/")  // http.../meta/copy/PATH/@activity/@time:TIME

	// File paths
	filesTopLevelSuffix  = "files"
	filesCreateSuffix    = strings.Join([]string{filesTopLevelSuffix, "/%v?", overwriteOption}, "")                  // http.../files/PATH?overwrite=FLAG
	filesRenameSuffix    = strings.Join([]string{filesTopLevelSuffix, "/%v?", nameOption, "&", overwriteOption}, "") // http.../files/PATH?name=NEWFILENAME&overwrite=FLAG
	filesMoveSuffix      = strings.Join([]string{filesTopLevelSuffix, "/%v?", pathOption, "&", overwriteOption}, "") // http.../files/PATH?overwrite=FLAG
	filesThumbnailSuffix = strings.Join([]string{filesTopLevelSuffix, "/%v?", sizeOption}, "")                       // http.../files/PATH?size=SIZE
)

func NewFileService(client *Client) *FileService {
	fs := new(FileService)
	fs.client = client
	return fs
}

// Returns the top level metadata (this is root folder, cannot change, see docs)
//
// https://www.copy.com/developer/documentation#api-calls/filesystem
func (fs *FileService) GetTopLevelMeta() (*Meta, error) {
	meta := new(Meta)
	resp, err := fs.client.DoRequestDecoding("GET", metaTopLevelSuffix, nil, meta)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 { // 400s and 500s
		return nil, errors.New(fmt.Sprintf("Client response: %d", resp.StatusCode))
	}

	return meta, nil
}

// Returns the metadata of a file
//
// https://www.copy.com/developer/documentation#api-calls/filesystem
func (fs *FileService) GetMeta(path string) (*Meta, error) {

	path = strings.Trim(path, "/")

	meta := new(Meta)
	resp, err := fs.client.DoRequestDecoding("GET", fmt.Sprintf(getMetaSuffix, path), nil, meta)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 { // 400s and 500s
		return nil, errors.New(fmt.Sprintf("Client response: %d", resp.StatusCode))
	}

	return meta, nil
}

// Returns all the metadata revisions of a file
//
// https://www.copy.com/developer/documentation#api-calls/filesystem
func (fs *FileService) ListRevisionsMeta(path string) ([]Revision, error) {
	meta := new(Meta)
	resp, err := fs.client.DoRequestDecoding("GET", fmt.Sprintf(listRevisionsSuffix, path), nil, meta)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 { // 400s and 500s
		return nil, errors.New(fmt.Sprintf("Client response: %d", resp.StatusCode))
	}

	return meta.Revisions, nil
}

// Returns the metadata in an specified revision
//
// https://www.copy.com/developer/documentation#api-calls/filesystem
func (fs *FileService) GetRevisionMeta(path string, time int) (*Meta, error) {
	return nil, nil
}

// Returns the file content. the user NEEDS TO CLOSE the buffer after using it
//
// https://www.copy.com/developer/documentation#api-calls/filesystem
func (fs *FileService) GetFile(path string) (io.ReadCloser, error) {

	resp, err := fs.client.DoRequestContent(strings.Join([]string{filesTopLevelSuffix, path}, "/"))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 { // 400s and 500s
		return nil, errors.New(fmt.Sprintf("Client response: %d", resp.StatusCode))
	}

	return resp.Body, nil
}

// Deletes the file content
//
// https://www.copy.com/developer/documentation#api-calls/filesystem
func (fs *FileService) DeleteFile(path string) error {
	return nil
}

// Uploads the file. Loads the file from the file path and uploads to the
// uploadPath.
// For example:
//   filePath: /home/slok/myFile.txt
//   UploadPath: test/uploads/something.txt
//
// https://www.copy.com/developer/documentation#api-calls/filesystem
func (fs *FileService) UploadFile(filePath, uploadPath string, overwrite bool) error {

	// Sanitize path
	uploadPath = strings.Trim(uploadPath, "/")

	// Get upload filename
	filename := filepath.Base(uploadPath)

	if filename == "" {
		return errors.New("Wrong uploadPath")
	}

	// Get upload path
	uploadPath = filepath.Dir(uploadPath)

	if uploadPath == "." { // Check if is at root, if so delete the point returned by Dir
		uploadPath = ""
	}

	// Sanitize path again
	uploadPath = strings.Trim(uploadPath, "/")

	// Create final path
	uploadPath = fmt.Sprintf(filesCreateSuffix, uploadPath, overwrite)

	res, err := fs.client.DoRequestMultipart(filePath, uploadPath, filename)

	if err != nil {
		return err
	}

	if res.StatusCode >= 400 { // 400s and 500s
		return errors.New(fmt.Sprintf("Client response: %d", res.StatusCode))
	}

	return nil
}

// Renames the file
//
// https://www.copy.com/developer/documentation#api-calls/filesystem
func (fs *FileService) RenameFile(path string, newName string, overwrite bool) error {
	path = strings.Trim(path, "/")
	return fs.moveOrRenameFile(fmt.Sprintf(filesRenameSuffix, path, newName, overwrite))
}

// Moves the file
//
// https://www.copy.com/developer/documentation#api-calls/filesystem
func (fs *FileService) MoveFile(path string, newPath string, overwrite bool) error {
	return nil
}

// Move and rename calls are similar, wrap in this function for convienence
func (fs *FileService) moveOrRenameFile(finalUrl string) error {

	resp, err := fs.client.DoRequestDecoding("PUT", finalUrl, nil, nil)

	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 { // 400s and 500s
		return errors.New(fmt.Sprintf("Client response: %d", resp.StatusCode))
	}

	return nil
}

// Creates a directory
//
// https://www.copy.com/developer/documentation#api-calls/filesystem
func (fs *FileService) CreateDirectory(path string) error {
	return nil
}
