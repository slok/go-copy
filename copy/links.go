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

func (linkService *LinkService) GetLink(token string) (*Meta, error) {
	return nil, nil
}

func (linkService *LinkService) GetLinks() ([]Meta, error) {
	return nil, nil
}

func (linkService *LinkService) CreateLink(name string, paths []string, public bool) (*Meta, error) {
	return nil, nil
}

func (linkService *LinkService) AddPaths(token string, paths []string) (*Meta, error) {
	return nil, nil
}

func (linkService *LinkService) AddRecipients(token string, recipients []Recipient) (*Meta, error) {
	return nil, nil
}

func (linkService *LinkService) DeleteLink(token string) (*Meta, error) {
	return nil, nil
}

func (linkService *LinkService) GetFilesMetaFromLink(token string) (*Meta, error) {
	return nil, nil
}
