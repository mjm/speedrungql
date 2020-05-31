package speedrungql

import (
	"context"
)

func (c *Client) GetRun(ctx context.Context, runID string) (*Run, error) {
	path := "/runs/" + runID
	var run Run
	if err := c.loadItem(ctx, path, &run); err != nil {
		return nil, err
	}
	return &run, nil
}
