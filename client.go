package speedrungql

import (
	"net/http"

	"github.com/graph-gophers/dataloader"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string

	loader *dataloader.Loader
}

func NewClient(baseURL string) *Client {
	c := &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
	c.loader = c.newLoader()
	return c
}
