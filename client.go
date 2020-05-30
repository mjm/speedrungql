package speedrungql

import (
	"net/http"

	"github.com/graph-gophers/dataloader"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string

	platformLoader *dataloader.Loader
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
	c.platformLoader = c.newLoader(func(key dataloader.Key) string {
		return "/platforms/" + key.String()
	})
}

func keysFromIDs(ids []string) dataloader.Keys {
	var keys dataloader.Keys
	for _, id := range ids {
		keys = append(keys, dataloader.StringKey(id))
	}
	return keys
}
