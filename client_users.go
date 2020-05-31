package speedrungql

import (
	"context"
)

func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	path := "/users/" + userID
	var user User
	if err := c.loadItem(ctx, path, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
