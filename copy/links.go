package copy

import (
	"strings"
)

// Use the object Meta from files (move to other file?)

type LinkService struct {
	client *Client
}

var (
	// Links paths
	linksTopLevelSuffix = "links"
	linksGetSuffix      = strings.Join([]string{linksTopLevelSuffix, "%v"}, "/") // https://.../links/TOKEN
)

func NewLinkService(client *Client) *LinkService {
	fs := new(LinkService)
	fs.client = client
	return fs
}

func (ls *LinkService) GetLink(token string) (*Meta, error) {
	return nil, nil
}

func (ls *LinkService) GetLinks() ([]Meta, error) {
	return nil, nil
}

func (ls *LinkService) CreateLink(name string, paths []string, public bool) error {
	return nil
}

func (ls *LinkService) AddPaths(token string, paths []string) error {
	return nil
}

func (ls *LinkService) AddRecipients(token string, recipients []Recipient) error {
	return nil
}

func (ls *LinkService) DeleteLink(token string) error {
	return nil
}

func (ls *LinkService) GetFilesMetaFromLink(token string) (*Meta, error) {
	return nil, nil
}
