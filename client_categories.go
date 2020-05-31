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

func (c *Client) GetCategory(ctx context.Context, categoryID string) (*Category, error) {
	path := "/categories/" + categoryID
	var category Category
	if err := c.loadItem(ctx, path, &category); err != nil {
		return nil, err
	}
	return &category, nil
}
