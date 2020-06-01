package speedrun

import (
	"context"
	"fmt"
)

func (c *Client) GetLeaderboard(ctx context.Context, gameID string, categoryID string, levelID *string, opts ...FetchOption) (*Leaderboard, error) {
	var path string
	if levelID == nil {
		path = fmt.Sprintf("/leaderboards/%s/category/%s", gameID, categoryID)
	} else {
		path = fmt.Sprintf("/leaderboards/%s/level/%s/category/%s", gameID, *levelID, categoryID)
	}

	var resp LeaderboardResponse
	if err := c.fetch(ctx, path, &resp, opts...); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
