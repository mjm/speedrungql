package speedrungql

import (
	"context"
	"fmt"
	"strings"
)

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
