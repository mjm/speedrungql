package speedrun

import (
	"context"
	"fmt"
	"strings"
)

func (c *Client) ListUsers(ctx context.Context, opts ...FetchOption) ([]*User, *PageInfo, error) {
	var resp UsersResponse
	if err := c.fetch(ctx, "/users", &resp, opts...); err != nil {
		return nil, nil, err
	}
	return resp.Data, resp.Pagination, nil
}

func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	var user User
	if err := c.loadItem(ctx, c.userKey(userID), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) userKey(id string) string {
	if strings.HasPrefix(id, c.BaseURL) {
		return id
	}
	return fmt.Sprintf("%s/users/%s", c.BaseURL, id)
}
