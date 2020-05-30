package speedrungql

import (
	"net/http"

	"github.com/graph-gophers/dataloader"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string

	platformLoader *dataloader.Loader
	userLoader     *dataloader.Loader
}

func NewClient(baseURL string) *Client {
	c := &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
	c.createLoaders()
	return c
}

func (c *Client) createLoaders() {
	c.platformLoader = c.newPlatformLoader()
	c.userLoader = c.newUserLoader()
}

func keysFromIDs(ids []string) dataloader.Keys {
	var keys dataloader.Keys
	for _, id := range ids {
		keys = append(keys, dataloader.StringKey(id))
	}
	return keys
}
