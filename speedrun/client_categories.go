package speedrun

import (
	"context"
	"fmt"
	"strings"
)

func (c *Client) ListGameCategories(ctx context.Context, gameID string, opts ...FetchOption) ([]*Category, error) {
	var resp CategoriesResponse
	if err := c.fetch(ctx, fmt.Sprintf("/games/%s/categories", gameID), &resp, opts...); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (c *Client) ListLevelCategories(ctx context.Context, levelID string, opts ...FetchOption) ([]*Category, error) {
	var resp CategoriesResponse
	if err := c.fetch(ctx, fmt.Sprintf("/levels/%s/categories", levelID), &resp, opts...); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (c *Client) GetCategory(ctx context.Context, categoryID string) (*Category, error) {
	var category Category
	if err := c.loadItem(ctx, c.categoryKey(categoryID), &category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *Client) categoryKey(id string) string {
	if strings.HasPrefix(id, c.BaseURL) {
		return id
	}
	return fmt.Sprintf("%s/categories/%s", c.BaseURL, id)
}
