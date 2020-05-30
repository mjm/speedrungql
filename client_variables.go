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
