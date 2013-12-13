package copy

import (
	"testing"
)

var (
	linkService *LinkService
)

func setupLinkService(t *testing.T) {
	setup(t)
	linkService = NewLinkService(client)
}

func tearDownLinkService() {
	defer tearDown()
}
