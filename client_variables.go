package speedrungql

import (
	"context"
	"fmt"
	"strings"
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

func (c *Client) ListLevelVariables(ctx context.Context, levelID string) ([]*Variable, error) {
	var resp VariablesResponse
	if err := c.fetch(ctx, fmt.Sprintf("/levels/%s/variables", levelID), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) GetVariable(ctx context.Context, variableID string) (*Variable, error) {
	var v Variable
	if err := c.loadItem(ctx, c.variableKey(variableID), &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (c *Client) variableKey(id string) string {
	if strings.HasPrefix(id, c.BaseURL) {
		return id
	}
	return fmt.Sprintf("%s/variables/%s", c.BaseURL, id)
}
