package speedrungql

import (
	"context"
	"encoding/json"
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

func (c *Client) GetPlatform(ctx context.Context, id string) (*Platform, error) {
	var platform Platform
	if err := c.loadItem(ctx, c.platformLoader, id, &platform); err != nil {
		return nil, err
	}
	return &platform, nil
}

func (c *Client) GetPlatforms(ctx context.Context, ids []string) ([]*Platform, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	ress, errs := c.platformLoader.LoadMany(ctx, keysFromIDs(ids))()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	var platforms []*Platform
	for _, res := range ress {
		var platform Platform
		if err := json.Unmarshal(res.(*EnvelopeResponse).Data, &platform); err != nil {
			return nil, err
		}
		platforms = append(platforms, &platform)
	}
	return platforms, nil
}

func keysFromIDs(ids []string) dataloader.Keys {
	var keys dataloader.Keys
	for _, id := range ids {
		keys = append(keys, dataloader.StringKey(id))
	}
	return keys
}
