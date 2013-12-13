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
