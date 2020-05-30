package speedrungql

import (
	"context"

	"github.com/graph-gophers/dataloader"
)

func (c *Client) newRunLoader() *dataloader.Loader {
	return c.newLoader(func(key dataloader.Key) string {
		return "/runs/" + key.String()
	})
}

func (c *Client) GetRun(ctx context.Context, runID string) (*Run, error) {
	var run Run
	if err := c.loadItem(ctx, c.runLoader, runID, &run); err != nil {
		return nil, err
	}
	return &run, nil
}
