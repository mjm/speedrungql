package speedrungql

import (
	"context"

	"github.com/graph-gophers/dataloader"
)

func (c *Client) newGameLoader() *dataloader.Loader {
	return c.newLoader(func(key dataloader.Key) string {
		return "/games/" + key.String()
	})
}

func (c *Client) ListGames(ctx context.Context, opts ...FetchOption) ([]*Game, *PageInfo, error) {
	var resp GamesResponse
	if err := c.fetch(ctx, "/games", &resp, opts...); err != nil {
		return nil, nil, err
	}

	return resp.Data, resp.Pagination, nil
}

func (c *Client) GetGame(ctx context.Context, gameID string) (*Game, error) {
	var game Game
	if err := c.loadItem(ctx, c.gameLoader, gameID, &game); err != nil {
		return nil, err
	}
	return &game, nil
}
