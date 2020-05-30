package speedrungql

import (
	"context"
	"fmt"
)

func (c *Client) ListGameCategories(ctx context.Context, gameID string, opts ...FetchOption) ([]*Category, error) {
	var resp CategoriesResponse
	if err := c.fetch(ctx, fmt.Sprintf("/games/%s/categories", gameID), &resp, opts...); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
