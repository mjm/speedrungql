package speedrun

import (
	"context"
	"fmt"
	"strings"
)

func (c *Client) ListRuns(ctx context.Context, opts ...FetchOption) ([]*Run, *PageInfo, error) {
	var resp RunsResponse
	if err := c.fetch(ctx, "/runs", &resp, opts...); err != nil {
		return nil, nil, err
	}
	return resp.Data, resp.Pagination, nil
}

func (c *Client) GetRun(ctx context.Context, runID string) (*Run, error) {
	var run Run
	if err := c.loadItem(ctx, c.runKey(runID), &run); err != nil {
		return nil, err
	}
	return &run, nil
}

func (c *Client) runKey(id string) string {
	if strings.HasPrefix(id, c.BaseURL) {
		return id
	}
	return fmt.Sprintf("%s/runs/%s", c.BaseURL, id)
}
