package speedrungql

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/graph-gophers/dataloader"
)

func (c *Client) GetRegion(ctx context.Context, regionID string) (*Region, error) {
	var region Region
	if err := c.loadItem(ctx, c.regionKey(regionID), &region); err != nil {
		return nil, err
	}
	return &region, nil
}

func (c *Client) GetRegions(ctx context.Context, ids []string) ([]*Region, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var keys dataloader.Keys
	for _, id := range ids {
		keys = append(keys, dataloader.StringKey(c.regionKey(id)))
	}

	ress, errs := c.loader.LoadMany(ctx, keys)()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	var regions []*Region
	for _, res := range ress {
		var region Region
		if err := json.Unmarshal(res.(*EnvelopeResponse).Data, &region); err != nil {
			return nil, err
		}
		regions = append(regions, &region)
	}
	return regions, nil
}

func (c *Client) regionKey(id string) string {
	if strings.HasPrefix(id, c.BaseURL) {
		return id
	}
	return fmt.Sprintf("%s/regions/%s", c.BaseURL, id)
}
