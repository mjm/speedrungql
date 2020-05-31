package speedrungql

import (
	"context"
	"fmt"
	"strings"
)

func (c *Client) ListGameLevels(ctx context.Context, gameID string) ([]*Level, error) {
	var resp LevelsResponse
	if err := c.fetch(ctx, fmt.Sprintf("/games/%s/levels", gameID), &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (c *Client) GetLevel(ctx context.Context, levelID string) (*Level, error) {
	var level Level
	if err := c.loadItem(ctx, c.levelKey(levelID), &level); err != nil {
		return nil, err
	}
	return &level, nil
}

func (c *Client) levelKey(id string) string {
	if strings.HasPrefix(id, c.BaseURL) {
		return id
	}
	return fmt.Sprintf("%s/levels/%s", c.BaseURL, id)
}
