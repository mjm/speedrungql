package speedrun

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/graph-gophers/dataloader"
)

func (c *Client) GetEngine(ctx context.Context, engineID string) (*Engine, error) {
	var category Engine
	if err := c.loadItem(ctx, c.engineKey(engineID), &category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *Client) GetEngines(ctx context.Context, ids []string) ([]*Engine, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var keys dataloader.Keys
	for _, id := range ids {
		keys = append(keys, dataloader.StringKey(c.engineKey(id)))
	}

	ress, errs := c.loader.LoadMany(ctx, keys)()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	var engines []*Engine
	for _, res := range ress {
		var engine Engine
		if err := json.Unmarshal(res.(*EnvelopeResponse).Data, &engine); err != nil {
			return nil, err
		}
		engines = append(engines, &engine)
	}
	return engines, nil
}

func (c *Client) engineKey(id string) string {
	if strings.HasPrefix(id, c.BaseURL) {
		return id
	}
	return fmt.Sprintf("%s/engines/%s", c.BaseURL, id)
}
