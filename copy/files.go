package copy

import (
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
	Stub                    bool                    `json:"stub,omitempty"`
	Share                   bool                    `json:"share,omitempty"`
	Children                []Meta                  `json:"children,omitempty"` // Inception :D
	Counts                  Count                   `json:"counts,omitempty"`   // Array? (sometimes? ask copy.com)
	RecipientConfirmed      bool                    `json:""recipient_confirmed",omitempty"`
	ObjectAvailable         bool                    `json:"object_available,omitempty"`
	Links                   []Link                  `json:"links,omitempty"`
	Revisions               []Revision              `json:"revisions,omitempty"`
	Url                     string                  `json:"url,omitempty"`
	RevisionId              int                     `json:"revision_id,omitempty"`
	Thumb                   string                  `json:"thumb,omitempty"`
	ThumbOriginalDimensions ThumbOriginalDimensions `json:"thumb_original_dimensions,omitempty"`
	ChildrenCount           int                     `json:"children_count",omitempty"`
}

type Count struct {
	New    int `json:"new,omitempty"`
	Viewed int `json:"viewed,omitempty"`
	Hidden int `json:"hidden,omitempty"`
}

type Link struct {
	Id                    string      `json:"id,omitempty"`
	Public                bool        `json:"public,omitempty"`
	Expires               bool        `json:"expires,omitempty"`
	Expired               bool        `json:"expired,omitempty"`
	Url                   string      `json:"url,omitempty"`
	UrlShort              string      `json:"url_short,omitempty"`
	Recipients            []Recipient `json:"recipients,omitempty"`
	CreatorId             string      `json:"creator_id,omitempty"`
	Confirmation_required bool        `json:"confirmation_required,omitempty"`
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
	metaTopLevelSuffix  = "meta"
	firstLevelSuffix    = strings.Join([]string{metaTopLevelSuffix, "copy"}, "/")
	filesTopLevelSuffix = "files"
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
	fs.client.Do("GET", metaTopLevelSuffix, nil, meta)
	return meta, nil
}
