package speedrungql

import (
	"context"
	"fmt"
)

func (c *Client) ListGameVariables(ctx context.Context, gameID string) ([]*Variable, error) {
	var resp VariablesResponse
	if err := c.fetch(ctx, fmt.Sprintf("/games/%s/variables", gameID), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) ListCategoryVariables(ctx context.Context, categoryID string) ([]*Variable, error) {
	var resp VariablesResponse
	if err := c.fetch(ctx, fmt.Sprintf("/categories/%s/variables", categoryID), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
