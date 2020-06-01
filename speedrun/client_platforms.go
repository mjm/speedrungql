package speedrun

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/graph-gophers/dataloader"
)

func (c *Client) ListPlatforms(ctx context.Context, opts ...FetchOption) ([]*Platform, *PageInfo, error) {
	var resp PlatformsResponse
	if err := c.fetch(ctx, "/platforms", &resp, opts...); err != nil {
		return nil, nil, err
	}

	return resp.Data, resp.Pagination, nil
}

func (c *Client) GetPlatform(ctx context.Context, id string) (*Platform, error) {
	var platform Platform
	if err := c.loadItem(ctx, c.platformKey(id), &platform); err != nil {
		return nil, err
	}
	return &platform, nil
}

func (c *Client) GetPlatforms(ctx context.Context, ids []string) ([]*Platform, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var keys dataloader.Keys
	for _, id := range ids {
		keys = append(keys, dataloader.StringKey(c.platformKey(id)))
	}

	ress, errs := c.loader.LoadMany(ctx, keys)()
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

func (c *Client) platformKey(id string) string {
	if strings.HasPrefix(id, c.BaseURL) {
		return id
	}
	return fmt.Sprintf("%s/platforms/%s", c.BaseURL, id)
}
