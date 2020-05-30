package speedrungql

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader"
)

func (c *Client) newCategoryLoader() *dataloader.Loader {
	return c.newLoader(func(key dataloader.Key) string {
		return "/categories/" + key.String()
	})
}

func (c *Client) ListGameCategories(ctx context.Context, gameID string, opts ...FetchOption) ([]*Category, error) {
	var resp CategoriesResponse
	if err := c.fetch(ctx, fmt.Sprintf("/games/%s/categories", gameID), &resp, opts...); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (c *Client) GetCategory(ctx context.Context, categoryID string) (*Category, error) {
	var category Category
	if err := c.loadItem(ctx, c.categoryLoader, categoryID, &category); err != nil {
		return nil, err
	}
	return &category, nil
}
