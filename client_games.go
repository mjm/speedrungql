package speedrungql

import (
	"context"
)

func (c *Client) ListGames(ctx context.Context, opts ...FetchOption) ([]*Game, *PageInfo, error) {
	var resp GamesResponse
	if err := c.fetch(ctx, "/games", &resp, opts...); err != nil {
		return nil, nil, err
	}

	return resp.Data, resp.Pagination, nil
}
