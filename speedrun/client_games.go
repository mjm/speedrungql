package speedrun

import (
	"context"
	"fmt"
	"strings"
)

func (c *Client) ListGames(ctx context.Context, opts ...FetchOption) ([]*Game, *PageInfo, error) {
	var resp GamesResponse
	if err := c.fetch(ctx, "/games", &resp, opts...); err != nil {
		return nil, nil, err
	}

	return resp.Data, resp.Pagination, nil
}

func (c *Client) GetGame(ctx context.Context, gameID string) (*Game, error) {
	var game Game
	if err := c.loadItem(ctx, c.gameKey(gameID), &game); err != nil {
		return nil, err
	}
	return &game, nil
}

func (c *Client) gameKey(id string) string {
	if strings.HasPrefix(id, c.BaseURL) {
		return id
	}
	return fmt.Sprintf("%s/games/%s", c.BaseURL, id)
}
